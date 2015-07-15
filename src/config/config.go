package config

import (
	"os"
	//	"os/exec"
	//"path/filepath"

	. "fmt"
)

const (
	DB_NAME     string = "jiamin1"
	DB_HOST     string = "i.mysqldev.mix.sina.com.cn"
	DB_PORT     string = "3306"
	DB_PASSWORD string = "root007"
	DB_USER     string = "root"
	//ROOT        string = "/home/jiamin1/go/src/"
	//VIEW        string = "/home/jiamin1/go/src/view/"

	//URL [2]string = {"baiu.com", "sina.cn"}
)

var ROOT string
var VIEW string

type Host struct {
	Url map[string]int
	//Url interface{}
}

//var Host Host

/*

//Host.
*/
func Init() Host {
	var res Host
	//func (Host *Host) Init() {
	tmp := map[string]int{
		"baiu.com": 1,
		"sina.cn":  2,
	}
	res.Url = tmp
	pwd, _ := os.Getwd()
	//动态更新root，view，防止不同项目，不同值
	ROOT = pwd + "/../"
	VIEW = ROOT + "view/"
	Println(VIEW)
	return res
}
