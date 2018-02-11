package util

import (
	"time"
)

type Tasks struct {
	invl     time.Duration
	InChan   chan []Task
	inputs   []Task
	procPriv func(interface{})
	outChan  chan []interface{}
}

type Task struct {
	Priv interface{}
	Ttl  int
}

func NewProc(invl time.Duration, procPriv func(interface{}), ret chan []interface{}) (fd *Tasks) {
	return &Tasks{
		invl:     invl,
		InChan:   make(chan []Task, 10),
		inputs:   make([]Task, 0),
		procPriv: procPriv,
		outChan:  ret,
	}
}

func (f *Tasks) Proc() {
	invl := time.NewTimer(f.invl).C

	for {
		select {
		case infos := <-f.InChan:
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
			f.procDetail()

			inputs := make([]Task, 0)
			for _, info := range f.inputs {
				tNow := int(time.Now().Unix())
				if info.Ttl > tNow {
					inputs = append(inputs, info)
				}
			}
			f.inputs = inputs

			invl = time.NewTimer(f.invl).C
		}
	}
}

func (f *Tasks) procDetail() {
	outputs := make([]interface{}, 0)
	for _, info := range f.inputs {
		f.procPriv(info.Priv)
		outputs = append(outputs, info.Priv)
	}

	f.outChan <- outputs
}
