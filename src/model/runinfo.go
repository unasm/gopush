package model

import (
	. "fmt"
	"runtime/pprof"
	"time"
)

func Runstate() {
	for {
		profiles := pprof.Profiles()
		for k, v := range profiles {
			Println("key is :", k)
			Println("value is :", v)
		}
		time.Sleep(time.Second * 10)
		Println()
		Println()
	}
}
