package inspect

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
	for {
		Runstate()
		CpuState()
		NetState()
		time.Sleep(time.Second * 100)
	}
}

//从golang内部 查看内存的使用情况
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
	显示网络的状态
	 Ipkts/s   The number of packets received per second.

	Ibytes/s  The number of bytes received per second.

	Opkts/s   The number of packets sent per second.

	Obytes/s  The number of bytes sent per second.
*/
func NetState() {
	net := GetNetWork()
	for k, v := range net {
		Printf("device : %s\n", k)
		Printf("\t 每秒接收数据包: \t%d\n", v[0])
		Printf("\t 每秒接受:\t%d kb\n", v[1]/1024)
		Printf("\t 每秒发送数据包: \t%d\n", v[2])
		Printf("\t 每秒发送\t %d kB\n", v[3]/1024)

	}
}

/*
 *  IFACE 是设备名，
 */
func GetNetWork() map[string][]int64 {
	var buf []int64
	cmd := exec.Command("sar", "-n", "DEV", "1", "1")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	res := make(map[string][]int64)
	var index int
	for i := 0; true; i = 1 {
		line, err := out.ReadString('\n')
		if i == 0 {
			continue
		}
		if err != nil {
			break
		}
		tokens := strings.Split(line, " ")
		var cnt int64
		cnt = 0
		buf = make([]int64, len(line))
		index = 0
		for j := 0; j < len(tokens); j++ {
			tokens[j] = strings.TrimSpace(tokens[j])
			if tokens[j] != "" {
				tmp, err := strconv.ParseInt(tokens[j], 10, 64)
				if err == nil {
					buf[index] = tmp
					cnt += tmp
					index++
				}
			}
		}
		if index == 4 && cnt > 0 && tokens[0] == "Average:" {
			res[tokens[1]] = make([]int64, index)
			for k := 0; k < index; k++ {
				res[tokens[1]][k] = buf[k]
			}
		}
	}
	/*
		for k, line := range res {
			Printf("%s\n", k)
			for key, v := range line {
				Printf("%d, %d\t", key, v)
			}
			Printf("\n")
		}
	*/
	return res
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
