package model

import (
	"bytes"
	. "fmt"
	"log"
	//	"os"
	"os/exec"
	"runtime"
	//"runtime/pprof"
	"strconv"
	"strings"

	"time"
)

var m runtime.MemStats
var stack []runtime.StackRecord

const (
	ShowCpuDetail = false
)

func Inspect() {
	Runstate()
	CpuState()
	time.Sleep(time.Second * 100)
}

//查看内存的使用情况
func Runstate() {
	runtime.ReadMemStats(&m)
	Println("堆内存为 : ", m.HeapSys/1000, "MB")
	Println("闲置的堆内存 : ", m.HeapIdle/1000, "MB")
	Println("内存总量 : ", m.Sys/1000, " MB")
	Println("使用中的内存 : ", m.Alloc/1000, " MB")
	Println("内存分配次数 : ", m.Mallocs)
	Println("内存释放次数 : ", m.Frees)
	Println("分配的系统栈容量 : ", m.StackSys/1000, "MB\n")
	Println("目前开启的goroutine数量为 : ", runtime.NumGoroutine())
	Println("CPU的数量为 : ", runtime.NumCPU())
}

//如果发生了一场，导致goroutine意外退出
func CheckRecover() {
	Println("goroutine 意外退出")
}

//cpu的使用情况
func CpuState() {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	type Process struct {
		pid int
		cpu float64
	}
	if err != nil {
		log.Fatal(err)
	}
	var cpuCnt float64
	cpuCnt = 0
	for i := 0; true; i = 1 {
		line, err := out.ReadString('\n')
		if i == 0 {
			continue
		}
		if err != nil {
			break
		}
		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}

		cpu, err := strconv.ParseFloat(ft[2], 4)
		if err != nil {
			log.Fatal(err)
		}
		if cpu > 0 {
			//显示每个程序占用cpu情况
			if ShowCpuDetail {
				pid, err := strconv.Atoi(ft[1])
				if err != nil {
					continue
				}
				log.Println("Process ", pid, " takes ", cpu, " % of the CPU")
				Println(ft[len(ft)-1])
			}
			cpuCnt = cpu + cpuCnt
		}
	}
	Printf("now cpu cost is %.2f\n", cpuCnt)
}

/*
func Copy() {
	check := "lookup heap"
	switch check {
	case "lookup heap":
		p := pprof.Lookup("heap")
		p.WriteTo(os.Stdout, 2)
	case "lookup threadcreate":
		p := pprof.Lookup("threadcreate")
		p.WriteTo(os.Stdout, 2)
	case "lookup block":
		p := pprof.Lookup("block")
		p.WriteTo(os.Stdout, 2)
	case "start cpuprof":
		if cpuProfile == nil {
			if f, err := os.Create("game_server.cpuprof"); err != nil {
				log.Printf("start cpu profile failed: %v", err)
			} else {
				log.Print("start cpu profile")
				pprof.StartCPUProfile(f)
				cpuProfile = f
			}
		}
	case "stop cpuprof":
		if cpuProfile != nil {
			pprof.StopCPUProfile()
			cpuProfile.Close()
			cpuProfile = nil
			log.Print("stop cpu profile")
		}
	case "get memprof":
		if f, err := os.Create("game_server.memprof"); err != nil {
			log.Printf("record memory profile failed: %v", err)
		} else {
			runtime.GC()
			pprof.WriteHeapProfile(f)
			f.Close()
			log.Print("record memory profile")
		}
	}
}
*/
