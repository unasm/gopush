//应该分不同的项目，数组各不相同,refer不同，访问不同,不然管理机和前端机访问的url会造成冲突
package main

import (
	"config"
	. "fmt"
	"html/template"
	"log"
	"model"
	"model/wb"
	//"model/href"
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
	if model.IsExistOne("pid") == false {
		Println("No pid")
		return
	}
	if model.Isint("pid", r.Form["pid"][0]) == false {
		Println("aid is not int")
		return
	}
	pid, err := strconv.Atoi(r.Form["pid"][0])
	model.CheckErr(err, "转换失败")

	wb.Send("one more")
	wb.StartNew(pid, "", conn)
}

func Home(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("/home/jiamin1/test.html")
	Println(config.VIEW + "inspect.html")

	t, err := template.ParseFiles(config.VIEW + "inspect.html")
	//href.GetHost(r.URL)
	model.CheckErr(err, "解析模板失败")
	str := wb.Report()
	//str := "hello,world"
	type Data struct {
		Str template.HTML
	}
	t.Execute(w, Data{Str: template.HTML(str)})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//cc := config.Host{}
	config.Init()
	wb.Init()
	inspect.Inspect()
	defer func() {
		if r := recover(); r != nil {
			model.CheckRecover()
		}
	}()
	//go model.Copy()
	http.HandleFunc("/chat/", Index)
	http.HandleFunc("/inspect/", Home)
	//http.HandleFunc("/chat/", Home)
	//go model.Run()
	//go h.run()
	//wb.StartListenBroad()
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal("Listen and add serve error ", err)
	}
}
