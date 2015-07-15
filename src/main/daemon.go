//不死小强
package main

import (
	. "fmt"
	//"log"
	//"encoding/binary"
	"os"
	"os/exec"
	"strconv"
	"strings"
	//"syscall"
	"bytes"
	"time"
)

//开始监听的程序
func Start(cmdStr string) {
	Println(cmdStr)
	for {
		cmd := exec.Command(cmdStr)
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

func parseOutLine(out *bytes.Buffer) ([]string, bool) {
	var tokens, res []string
	line, err := out.ReadString('\n')
	if err != nil {
		return tokens, false
	}
	tokens = strings.Split(line, " ")
	res = make([]string, 0, len(tokens))
	for tmp, t := range tokens {
		if tmp == 0 {
			continue
		}
		if t != "" && t != "\t" {
			res = append(res, t)
			//Printf("%d : %s\t", tmp, t)
		}
	}
	//Printf("\n")
	return res, true
}

//执行系统的调用,然后格式化返回
func SysCall(cmdStr string, argc string) [][]string {
	var out bytes.Buffer
	res := make([][]string, 0, 100)
	cmd := exec.Command(cmdStr, argc)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		Println("exec run failed")
		return res
	}
	tmp, flag := parseOutLine(&out)
	if flag == false {
		return res
	}
	res = append(res, tmp)
	for {
		tmp, flag = parseOutLine(&out)
		if flag == false {
			return res
		}
		res = append(res, tmp)
	}
	return res
}
func main() {
	file, err := os.OpenFile("pid", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		os.Exit(1)
	}
	buf := make([]byte, 10)
	_, _ = file.Read(buf)
	//去除多余的0,不然会导致错误
	num := strings.Trim(string(buf), "\x00")
	if tmp, err := strconv.Atoi(num); err == nil && tmp > 0 {
		cmdAll := SysCall("ps", "aux")
		flag := 0
		for i, length := 1, len(cmdAll); i < length; i++ {
			if strings.EqualFold(cmdAll[i][0], num) {
				flag = 1
				break
			}
		}
		if flag == 1 {
			os.Exit(1)
		}
	} else if len(num) > 0 {
		if err != nil {
			Println("one wrong")
		}
		if tmp > 0 {
			Println("sec wrong")
		}
		os.Exit(1)
	}
	file.Seek(0, 0)
	defer file.Close()
	file.WriteString(Sprintf("%d", os.Getpid()))
	cmdStr := os.Getenv("GOPATH") + "/src/main/main"
	Start(cmdStr)
}
