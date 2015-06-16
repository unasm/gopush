<<<<<<< HEAD
//不死小强
=======
>>>>>>> 2e06b5628f2fdb626c4efbe3e01e9823951a6591
package main

import (
	. "fmt"
	//"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()
	//l := log.New(file, "", os.O_APPEND)

	for {
<<<<<<< HEAD
		cmd := exec.Command("/home/jiamin1/websocket")
=======
		cmd := exec.Command("/Users/tianyi/main")
>>>>>>> 2e06b5628f2fdb626c4efbe3e01e9823951a6591
		err := cmd.Start()
		if err != nil {
			Printf("%s 启动命令失败", time.Now(), err)
			time.Sleep(time.Second * 5)
			continue
		}
		//l.Printf("%s 进程启动", time.Now(), err)
		Println("进程启动 ", time.Now(), err)
		err = cmd.Wait()
		Println("进程退出 ", time.Now(), err)
		time.Sleep(time.Second * 5)
	}
}
