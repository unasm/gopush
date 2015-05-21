//go的并发编程
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func routine(name string, delay time.Duration) {
	t0 := time.Now()
	fmt.Println(name, "start at", t0)
	time.Sleep(delay)
	t1 := time.Now()
	fmt.Println(name, "end is ", t1)
	fmt.Println(name, "lasted ", t1.Sub(t0))
}
func main() {
	runtime.GOMAXPROCS(4)
	rand.Seed(time.Now().Unix())
	fmt.Println("time is %s", time.Now().Unix())
	var name string
	for i := 0; i < 3; i++ {
		name = fmt.Sprintf("go_%02d", i)
		go routine(name, time.Duration(rand.Intn(5))*time.Second)
	}
	var input string
	fmt.Scanln(&input)
	fmt.Println("done\n")
}
