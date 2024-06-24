package watcher

import (
	"fmt"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	watcher := NewWatcher()
	process := &Process{
		ID: "p1",
	}
	work := &Work{
		ID: "w1",
		FuncWork: func(work *Work) (interface{}, error) {
			return 1, nil
		},
	}
	watcher.AddProcess(process)
	process.AddWork(work)
	process.InitProcess()
	time.Sleep(time.Second * 1)
	fmt.Println(watcher.CheckProcessExists(process.ID))
	fmt.Println(process.Return)
}
