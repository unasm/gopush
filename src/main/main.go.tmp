package main

import (
	"config"
	. "fmt"
	"html/template"
	"log"
	"model"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	//	"runtime/pprof"
	"github.com/gorilla/websocket"
	//"reflect"
	"strconv"
	//"time"
)

/*
func unescaped(x string) interface{} { return template.HTML(x) }
func renderTemplate(w http.ResponseWriter, tmpl string, view *Page) {
	t := template.New("")
	t = t.Funcs(template.FuncMap{"unescaped": unescaped})
	t, err := t.ParseFiles("view.html", "edit.html")
	err = t.ExecuteTemplate(w, tmpl+".html", view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
*/
func Index(w http.ResponseWriter, r *http.Request) {

	if r.Proto != "HTTP/1.1" {
		Fprintf(w, "request proto is not permit")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			model.CheckRecover()
		}
	}()
	var upgrader = websocket.Upgrader{
		ReadBufferSize: 1024, WriteBufferSize: 1024,
		//HandshakeTimeout: HANDSHAKE,
		//握手时间超过3s超时
		//不再检查请求源
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	model.CheckErr(err, "升级失败")
	if err != nil {
		Fprintf(w, "upgrader faied")
		return
	}
	r.ParseForm()
	model.Form = r.Form
	if model.IsExistOne("jid") == false {
		Println("No jid")
		return
	}
	if model.Isint("jid", r.Form["jid"][0]) == false {
		Println("aid is not int")
		return
	}
	aid, err := strconv.Atoi(r.Form["jid"][0])
	model.CheckErr(err, "转换失败")
	if model.IsExistOne("editor") == false {
		Println("editor is NULL")
		return
	}
	uid := r.Form["editor"][0]
	ok, str := model.CheckRepeat(aid, uid)
	if ok == false {
		conn.WriteMessage(websocket.TextMessage, []byte(str))
		conn.Close()
		return
	}
	//Println(reflect.TypeOf(conn))
	model.StartNew(aid, uid, conn)
}

func Home(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("/home/jiamin1/test.html")
	t, err := template.ParseFiles(config.VIEW + "inspect.html")
	model.CheckErr(err, "解析模板失败")
	str := model.Report()
	type Data struct {
		Str template.HTML
	}
	t.Execute(w, Data{Str: template.HTML(str)})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//model.Connect()
	go model.Inspect()
	defer func() {
		if r := recover(); r != nil {
			model.CheckRecover()
		}
	}()
	//go model.Copy()
	http.HandleFunc("/chat/", Index)
	http.HandleFunc("/inspect/", Home)
	//http.HandleFunc("/chat/", Home)
	go model.Run()
	//go h.run()
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal("Listen and add serve error ", err)
	}
}
