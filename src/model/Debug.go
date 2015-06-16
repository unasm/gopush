//用于调试和输出错误信息
package model

import (
	. "fmt"
	"runtime"
)

//如果发生了一场，导致goroutine意外退出
func CheckRecover() {
	//Println("goroutine 意外退出")

}

// @param string	note	如果失败的话，需要说明的原因
// @param error		err		系统的报错
func CheckErr(err error, note string) {
	if err != nil {
		//将这里的信息记录在数据库里面
		Println(note)
		buf := make([]byte, 1<<20)
		runtime.Stack(buf, true)
		data := make(map[string]string)
		data["stack"] = string(buf)

		data["type"] = "exception"
		data["note"] = note
		data["err"] = err.Error()
		ErrorInsert(data)
		panic(err)
	}
}
