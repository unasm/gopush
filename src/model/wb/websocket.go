//定义websocket对象
package wb

import (
	. "fmt"
	"github.com/gorilla/websocket"
	"model"
	"time"
)

const (
	//每40s发送一次心跳信号，如果确认，延长关闭一份中
	PONGWAIT  = 4 * time.Second
	WRITEWAIT = 1 * time.Minute
	HANDSHAKE = 5 * time.Second
)

type Client struct {
	ws       *websocket.Conn
	shutdown chan bool
	art_id   int
	editor   string
}
type hub struct {
	Clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
}

var h = hub{
	//+1, 取消缓冲
	broadcast:  make(chan string, 1),
	register:   make(chan *Client, 1),
	unregister: make(chan *Client, 1),
	Clients:    make(map[*Client]bool),
}

//检查登陆者访问的连接是不是已经被访问了
func CheckRepeat(aid int, uid string) (bool, string) {
	var str string
	for link := range h.Clients {
		if aid == link.art_id {
			//回应一个已经加锁了
			if uid == link.editor {
				str = "您已经在本地打开了，请勿多处编辑"
				//conn.WriteMessage(websocket.TextMessage, []byte("您已经在本地打开了，请勿多处编辑"))
			} else {
				str = link.editor + "正在编辑，请稍等"
				//conn.WriteMessage(websocket.TextMessage, []byte(uid+"正在编辑，请稍等"))
			}
			//conn.Close()
			return false, str
		}
	}
	return true, ""
}

//向外广播推送内容
//阻塞监听broadcast 变化
func StartListenBroad() {
	for {
		select {
		case msg := <-h.broadcast:
			//可以根据数量，开多进程跑
			for conn, ok := range h.Clients {
				if ok {
					conn.ws.WriteMessage(websocket.TextMessage, []byte(msg))
				}
			}
		}
	}
}

//初始化wb中的变量，启动hub，
func Init() {
	go StartListenBroad()
	go ListenHub()
}

//向广播中发送数据
func Send(msg string) {
	h.broadcast <- msg
}

// 开始一个新的连接,将连接放在register中
func StartNew(aid int, uid string, conn *websocket.Conn) {
	c := &Client{
		shutdown: make(chan bool),
		ws:       conn,
		art_id:   aid,
		editor:   uid,
	}
	h.register <- c
	//Report()
	go c.Keeplink()
	c.readPump()
}

//回报目前的编辑情况
//func (h *hub)
func Report() string {
	str := ""
	//str += "<li>全部在线编辑 : " + len(h.Clients) + "</li>"
	str += Sprintf("<li> 全部在线编辑 %d </li>", len(h.Clients))
	for c, ok := range h.Clients {
		if ok {
			str += Sprintf("<li> %s is editing %d </li>", c.editor, c.art_id)
			//Println(c.editor, " is editing ", c.art_id)
		} else {
			str += "<li>false is existing</li>"
			//Println("false is existing")
		}
	}
	return str
}

//开始监听 hub中的变化
func ListenHub() {
	for {
		select {
		case Client := <-h.register:

			//Println("adding : ", Client.art_id)
			h.Clients[Client] = true
			break
		case Client := <-h.unregister:
			Println("before delete ", len(h.Clients), " : ", h.Clients[Client])
			delete(h.Clients, Client)
			Println("after delete ", len(h.Clients), " : ", h.Clients[Client])
			/*
				if h.Clients[Client] {
				}
			*/
			break
		}
	}
}

//保持心跳，维持通信
func (c *Client) Keeplink() {
	ticker := time.NewTicker(PONGWAIT)
	c.ws.SetReadDeadline(time.Now().Add(WRITEWAIT))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(WRITEWAIT))
		return nil
	})
	defer func() {
		ticker.Stop()
		if r := recover(); r != nil {
			model.CheckRecover()
		}
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

func (c *Client) readPump() {
	for {
		//其是可以用来向其他的连接推送广播的
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		//h.broadcast <- string(msg)
	}
	defer func() {
		if r := recover(); r != nil {
			model.CheckRecover()
		}
		c.ws.Close()
		h.unregister <- c
		c.shutdown <- true
	}()
}
