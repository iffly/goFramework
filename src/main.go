package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
	"util"
)

func main() {
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)

	retChan := make(chan []interface{})
	f := util.NewProc(10*time.Second, proc, retChan)
	go f.Proc()

	go input(f.InChan)
	go output(retChan)

	time.Sleep(120 * time.Second)
}

type infoPriv struct {
	name string
}

func proc(info interface{}) {
	priv, ok := info.(infoPriv)
	if !ok {
		fmt.Println("sth err")
		return
	}

	fmt.Println("proc priv", priv.name)
}

func input(inChan chan []util.Task) {
	j := 0

	for {
		// i := j
		i := 0

		for ; i < j+2; i++ {
			inChan <- []util.Task{util.Task{Priv: infoPriv{name: "infoPrivTest_" + strconv.Itoa(i)}, Ttl: int(time.Now().Unix()) + 20}}
			fmt.Println("input", i)
		}
		<-time.NewTimer(40 * time.Second).C

		// j = i
	}
}

func output(outChan chan []interface{}) {
	for {
		select {
		case infos := <-outChan:
			for _, info := range infos {
				one, ok := info.(infoPriv)
				if !ok {
					fmt.Println("sth err")
					continue
				}

				fmt.Println("output", one)
			}

		}
	}
}
