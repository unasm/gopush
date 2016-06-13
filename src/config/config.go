package config

import (
	"os"
	//	"os/exec"
	//"path/filepath"

	. "fmt"
)

const (
	DB_NAME     string = "qiye"
	DB_HOST     string = "100.73.17.75"
	DB_PORT     string = "3306"
	DB_PASSWORD string = "VMF9>4z@426Y"
	DB_USER     string = "devmysql"
)

var ROOT string
var VIEW string

type Host struct {
	Url map[string]int
}

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
