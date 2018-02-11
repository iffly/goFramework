package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func main() {
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)

	retChan := make(chan []Task)
	f := NewProc(10*time.Second, retChan)
	go f.Proc()

	go input(f.inChan)
	go output(retChan)

	time.Sleep(120 * time.Second)
}

func input(inChan chan []Task) {
	j := 0

	for {
		// i := j
		i := 0

		for ; i < j+2; i++ {
			inChan <- []Task{Task{name: "test_" + strconv.Itoa(i), ttl: int(time.Now().Unix()) + 20}}
			fmt.Println("input", i)
		}
		<-time.NewTimer(40 * time.Second).C

		// j = i
	}
}

func output(outChan chan []Task) {
	for {
		select {
		case infos := <-outChan:
			fmt.Println("output", infos)
		}
	}
}

type Tasks struct {
	inChan  chan []Task
	outChan chan []Task
	invl    time.Duration
	inputs  []Task
}

type Task struct {
	name string
	ttl  int
}

func NewProc(invl time.Duration, ret chan []Task) (fd *Tasks) {
	return &Tasks{
		inChan:  make(chan []Task, 10),
		outChan: ret,
		invl:    invl,
		inputs:  make([]Task, 0),
	}
}

func (f *Tasks) Proc() {
	invl := time.NewTimer(f.invl).C

	for {
		select {
		case infos := <-f.inChan:
			mark := false
			for _, new := range infos {
				for _, old := range f.inputs {
					if new == old {
						mark = true
						break
					}
				}
				if false == mark {
					f.inputs = append(f.inputs, new)
				}
				mark = false
			}
		case <-invl:
			f.ProcDetail()
			inputs := make([]Task, 0)
			for _, info := range f.inputs {
				tNow := int(time.Now().Unix())
				if info.ttl > tNow {
					inputs = append(inputs, info)
				}
			}
			f.inputs = inputs
			invl = time.NewTimer(f.invl).C
		}
	}
}

func (f *Tasks) ProcDetail() {
	outputs := make([]Task, 0)
	for _, info := range f.inputs {
		fmt.Println("proc", time.Now().Unix(), info.name)
		outputs = append(outputs, info)
	}
	f.outChan <- outputs
}
