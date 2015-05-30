package main

import (
	. "fmt"
	"github.com/gorilla/websocket"
	"log"
	"model"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const (
	//每40s发送一次心跳信号，如果确认，延长关闭一份中
	PONGWAIT  = 4 * time.Second
	WRITEWAIT = 1 * time.Minute
	HANDSHAKE = 5 * time.Second
)

type client struct {
	ws       *websocket.Conn
	shutdown chan bool
	art_id   int
	editor   string
}
type hub struct {
	clients    map[*client]bool
	broadcast  chan string
	register   chan *client
	unregister chan *client
}

var h = hub{
	//+1, 取消缓冲
	broadcast:  make(chan string, 1),
	register:   make(chan *client, 1),
	unregister: make(chan *client, 1),
	clients:    make(map[*client]bool),
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("panicing \n")
		panic(err)
	}
}
func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			//Println("adding : ", client.art_id)
			h.clients[client] = true
			break
		case client := <-h.unregister:
			//Println("cuting : ", client.art_id)
			if h.clients[client] {
				delete(h.clients, client)
			}
			break
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize: 1024, WriteBufferSize: 1024, HandshakeTimeout: HANDSHAKE,
		//握手时间超过3s超时
		//不再检查请求源
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	checkErr(err)
	r.ParseForm()
	model.Form = r.Form
	if model.IsExistOne("aid") == false {
		Println("No aid")
		return
	}
	if model.Isint("aid", r.Form["aid"][0]) == false {
		Println("aid is not int")
		return
	}
	aid, err := strconv.Atoi(r.Form["aid"][0])
	checkErr(err)
	if model.IsExistOne("uid") == false {
		Println("aid is NULL")
		return
	}
	uid := r.Form["uid"][0]
	for link := range h.clients {
		Println(link.art_id)
		if aid == link.art_id {
			//回应一个已经加锁了
			conn.WriteMessage(websocket.TextMessage, []byte(uid+" is editing"))
			conn.Close()
			return
		}
	}

	c := &client{
		shutdown: make(chan bool),
		ws:       conn,
		art_id:   aid,
		editor:   uid,
	}
	h.register <- c
	go c.Keeplink()
	c.readPump()
}

//保持心跳，维持通信
func (c *client) Keeplink() {
	ticker := time.NewTicker(PONGWAIT)
	c.ws.SetReadDeadline(time.Now().Add(WRITEWAIT))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(WRITEWAIT))
		return nil
	})
	defer func() {
		ticker.Stop()
		close(c.shutdown)
	}()
	for {
		select {
		case <-ticker.C:
			//发送一个ping信号,维持通信
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case flag := <-c.shutdown:
			if flag == true {
				Println("shuting down")
				return
			}
		}
	}
}

/**
 * 监听信道变化,捕捉close信号
 */

func (c *client) readPump() {
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		//h.broadcast <- string(msg)
	}
	defer func() {
		c.ws.Close()
		h.unregister <- c
		c.shutdown <- true
	}()

}
func main() {
	runtime.GOMAXPROCS(4)
	go model.Runstate()
	http.HandleFunc("/", Index)
	go h.run()
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal("Listen and add serve error ", err)
	}
}
