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

var conn_id int

/*
func Echo(ws *websocket.Conn) {
	var err error
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("cant receive")
			break
		}
		fmt.Println("Received back from client: " + reply)
		msg := "Received back from client: " + reply
		fmt.Println("send to client : " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

*/
func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("yes ,welcome !")
	fmt.Fprintf(w, "Hello,world<br/>")
}

//一种websocket的方式
func socket(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
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
func socket2(w http.ResponseWriter, r *http.Request) {
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
	conn.SetReadDeadline(time.Now().Add(60000000000))
	fmt.Println(conn_id)
	if err != nil {
		log.Println("not websocket: ", err)
		return
	}
	conn_id++
	for {
		fmt.Println(conn_id)
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
}

//路由处理
func route() {
	http.HandleFunc("/chat/", socket2)
	//http.HandleFunc("/chat/", socket)
	http.HandleFunc("/", index)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	route()
	//var ans *model.Result
	if err := http.ListenAndServe(":8010", nil); err != nil {
		conn_id = 0
		log.Fatal("Listen and Server : ", err)
	}
}
