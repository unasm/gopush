package main

import (
	"fmt"
	//"model"
	"net/http"
	//	"reflect"
	//"code.google.com/p/go.net/websocket"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"time"
)

//每刷新一次，就会建立一个连接，必须要根据用户的ip等，限制连接的创建
var conn_id int

const (
	//每40s发送一次心跳信号，如果确认，延长关闭一份中
	PONGWAIT  = 40 * time.Second
	WRITEWAIT = 1 * time.Minute
	HANDSHAKE = 5 * time.Second
)

type hub struct {
	clients    map[*client]bool
	broadcast  chan string
	register   chan *client
	unregister chan *client
}

type client struct {
	ws   *websocket.Conn
	send chan []byte
}

var h = hub{
	broadcast:  make(chan string),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
}

func (h *hub) broadcastMessage(content []byte) {
	fmt.Println(content)
	for c := range h.clients {
		select {
		//如果这里for循环的话，是线性的，如果通过信道的话，是并行的
		case c.send <- content:
			fmt.Println(content)
			break
		//case
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}

//扮演了信息中心的角色
func (h *hub) run() {
	for {
		//监控hub的几个信道变化
		select {
		case c := <-h.register:
			//新添加的节点
			h.clients[c] = true
			break
		case c := <-h.unregister:
			//连接中断的连接
			_, ok := h.clients[c]
			if ok {
				delete(h.clients, c)
				close(c.send)
			}
			break
		case m := <-h.broadcast:
			h.broadcastMessage([]byte(m))
			break
		}
	}
}

//正常情况下的http请求以及相应
func Index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("yes ,welcome !")
	fmt.Fprintf(w, "Hello,world323<br/>")
}

//第一种websocket的方式
func serverWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Cookies())
	fmt.Println("above is cookie")
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//握手时间超过3s超时
		HandshakeTimeout: HANDSHAKE,
		//不再检查请求源
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//将连接从http升级成websocket协议
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("not websocket: ", err)
		return
	}
	c := &client{
		send: make(chan []byte, 1024),
		ws:   conn,
	}
	h.register <- c
	go c.writePump()
	c.readPump()
}

/**
 *  从连接中读取信息，然后广播出去
 * 只有创建websocket的时候才会进来
 */
func (c *client) readPump() {
	defer func() {
		//窗口关闭或者连接关闭的时候，会触发这个defer
		//fmt.Println("close the connection")
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(1024 * 1024)
	c.ws.SetReadDeadline(time.Now().Add(WRITEWAIT))
	c.ws.SetWriteDeadline(time.Now().Add(WRITEWAIT))
	//根据ping信号，每次将连接断开的时间延后
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetWriteDeadline(time.Now().Add(WRITEWAIT))
		c.ws.SetReadDeadline(time.Now().Add(WRITEWAIT))
		return nil
	})
	for {
		fmt.Println(r.Cookies())
		//守护着这个连接，从连接中读取信息，然后通过信道传递给broadcast ，进入传递的队列
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		h.broadcast <- string(msg)
	}
}

//  守护进程，等着想对应的链接写入数据
func (c *client) writePump() {
	ticker := time.NewTicker(PONGWAIT)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		//监控send的信号变化,并维持通信畅通
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

/**
 */
func (c *client) write(mt int, msg []byte) error {
	return c.ws.WriteMessage(mt, msg)
}

//第一种websocket的方式
func Socket(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(http.Cookie())
	fmt.Println(r.Cookies())
	fmt.Println("above is cookie")
	//fmt.Println(r.Cookie())
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//握手时间超过3s超时
		HandshakeTimeout: HANDSHAKE,
		//不再检查请求源
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//将连接从http升级成websocket协议
	conn, err := upgrader.Upgrade(w, r, nil)

	conn.SetWriteDeadline(time.Now().Add(PONGWAIT))
	//conn.SetReadDeadline(time.Now().Add(600000000000))
	//设置读写的过期时间
	/*
		 c.ws.SetPongHandler(func(string) error {
		         c.ws.SetReadDeadline(time.Now().Add(pongWait));
			         return nil
				     })
	*/
	conn.SetPongHandler(func(string) error {
		fmt.Println("pong")
		fmt.Println(time.Now())
		conn.SetReadDeadline(time.Now().Add(60000000000))
		return nil
	})
	conn.SetPingHandler(func(string) error {
		fmt.Println("ping")
		fmt.Println(time.Now())
		conn.SetReadDeadline(time.Now().Add(60000000000))
		return nil
	})
	if err != nil {
		log.Println("not websocket: ", err)
		return
	}
	for {
		//读取数据，
		msg, p, err := conn.ReadMessage()
		fmt.Println(msg)
		fmt.Println(p)
		fmt.Println(err)
		if err != nil {
			return
		}
		//返回数据，
		if err = conn.WriteMessage(msg, p); err != nil {
			return
		}
	}
}

//第二种读取socket信息的方法
func Socket_bak(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		// ReadBufferSize and WriteBufferSize specify I/O buffer sizes. If a buffer
		// size is zero, then a default value of 4096 is used. The I/O buffer sizes
		// do not limit the size of the messages that can be sent or received.
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		//握手的超时时间,单位是ms
		HandshakeTimeout: 2000,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//将连接从http升级成websocket协议
	conn, err := upgrader.Upgrade(w, r, nil)
	//精确到纳秒 ， add 是time对象的方法，time.Now()产生一个time对象,time.Now()是此时的时间time对象，Add增加1分钟
	conn.SetWriteDeadline(time.Now().Add(60000000000))
	//设置读写的过期时间
	conn.SetReadDeadline(time.Now().Add(60000000000))
	if err != nil {
		log.Println("not websocket: ", err)
		return
	}
	fmt.Println("socket ")
	conn_id++
	for {
		//读取数据，msg是类型，1,2 ，数据的类型，字符还是二进制
		//p 是内容，err 是错误
		msg, p, err := conn.NextReader()
		if err != nil {
			log.Println("next Reader ", err)
			return
		}
		//返回数据，

		www, err := conn.NextWriter(msg)
		if err != nil {
			log.Println("next Writer ", err)
			return
		}
		if _, err := io.Copy(www, p); err != nil {
			log.Println("copy ", err)
			return
		}

		if err := www.Close(); err != nil {
			return
		}
	}
	//关闭conn
	defer conn.Close()
}

//路由处理
func route() {
	//http.HandleFunc("/chat/", Socket_bak)
	http.HandleFunc("/chat/", serverWs)
	//http.HandleFunc("/chat/", socket)
	//对路由为/的注册index函数
	http.HandleFunc("/", Index)
}

//检查错误，输出错误
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	conn_id = 0
	route()
	go h.run()
	// 监听 8010端口

	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal("Listen and Server : ", err)
	}
}
