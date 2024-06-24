package watcher

import (
	"go.uber.org/zap"
	"sync"
)

func NewWatcher() *Watcher {
	return &Watcher{}
}

type Watcher struct {
	processes []*Process

	mutex sync.RWMutex
}

func (w *Watcher) AddProcess(process *Process) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	process.ChangeState(TODO)

	w.processes = append(w.processes, process)
}

func (w *Watcher) CheckProcessExists(processID string) (exists bool, process *Process) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	for ii := range w.processes {
		if w.processes[ii].ID == processID {
			w.processes[ii].mutex.RLock()
			defer w.processes[ii].mutex.RUnlock()

			exists = true
			process = w.processes[ii]
			return
		}
	}

	process = &Process{ID: processID}

	return
}

func (w *Watcher) RemoveProcess(processID string) {
	for ii := range w.processes {
		if w.processes[ii].ID == processID {
			w.processes[ii].mutex.Lock()
			defer w.processes[ii].mutex.Unlock()
			w.processes = append(w.processes[:ii], w.processes[ii+1:]...)
		}
	}
}

type Process struct {
	ID     string
	works  []*Work
	Return any
	Error  error

	HaveStatePercentage
}

func (p *Process) AddWork(workIn *Work) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	workIn.ChangeState(TODO)

	p.works = append(p.works, workIn)

}

func (p *Process) InitProcess() {
	go func() {
		p.ChangeState(InProgress)

		var (
			returnedValue interface{}
			returnedErr   error
		)
		for ii := range p.works {
			p.works[ii].ChangeState(InProgress)

			returnedValue, returnedErr = p.works[ii].FuncWork(p.works[ii])

			if returnedErr != nil {
				p.works[ii].ChangeState(Error)
			} else {
				p.works[ii].ChangeState(Done)
			}

			if p.works[ii].State == Error {
				p.SetReturn(returnedValue, returnedErr)
				zap.L().Error("Error on process " + p.ID + " " + returnedErr.Error())
				return
			}
		}

		p.SetReturn(returnedValue, returnedErr)

		// TODO: not work
		if processReturn, ok := returnedValue.(*ProcessReturn); ok {
			state, _ := p.GetStatePercentage()
			processReturn.ChangeState(state)
		}
	}()
}

func (p *Process) CheckWorkExists(workID string) (exists bool, state State, percentage float64) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for ii := range p.works {
		if p.works[ii].ID == workID {
			p.works[ii].mutex.RLock()
			defer p.works[ii].mutex.RUnlock()

			exists = true
			state = p.works[ii].State
			percentage = p.works[ii].Percentage
			return
		}
	}

	return
}

func (p *Process) AddPercentageValue(percentage float64, noLock ...bool) {
	if len(noLock) == 0 || !noLock[0] {
		p.mutex.Lock()
		defer p.mutex.Unlock()
	}

	p.Percentage = percentage
}

func (p *Process) ChangeState(newState State, noLock ...bool) {
	if len(noLock) == 0 || !noLock[0] {
		p.mutex.Lock()
		defer p.mutex.Unlock()
	}
	switch newState {
	case TODO:
		p.AddPercentageValue(0, true)
	case Error, Done:
		p.AddPercentageValue(100, true)
	}

	p.State = newState
}

func (p *Process) SetReturn(returnedValue interface{}, returnedErr error, noLock ...bool) {
	if len(noLock) == 0 || !noLock[0] {
		p.mutex.Lock()
		defer p.mutex.Unlock()
	}
	p.Return = returnedValue
	p.Error = returnedErr

	if returnedErr != nil {
		p.ChangeState(Error, true)
	} else {
		p.ChangeState(Done, true)
	}
}

func (p *Process) Finished() bool {
	actualState, _ := p.GetStatePercentage()
	return actualState == Done || actualState == Error
}

func (p *Process) GetStatePercentage() (state State, percentage float64) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.State, p.Percentage
}

type Work struct {
	ID       string
	FuncWork func(work *Work) (interface{}, error)

	HaveStatePercentage
}

func (w *Work) AddPercentageValue(percentage float64, noLock ...bool) {
	if len(noLock) == 0 || !noLock[0] {
		w.mutex.Lock()
		defer w.mutex.Unlock()
	}

	w.Percentage = percentage
}

func (w *Work) ChangeState(newState State, noLock ...bool) {
	if len(noLock) == 0 || !noLock[0] {
		w.mutex.Lock()
		defer w.mutex.Unlock()
	}

	switch newState {
	case TODO:
		w.AddPercentageValue(0, true)
	case Error, Done:
		w.AddPercentageValue(100, true)
	}

	w.State = newState
}

type State string

type ProcessReturn struct {
	HaveStatePercentage
}

type HaveStatePercentage struct {
	mutex sync.RWMutex

	State      State
	Percentage float64
}

func (h *HaveStatePercentage) AddPercentageValue(percentage float64, noLock ...bool) {
	if len(noLock) == 0 || !noLock[0] {
		h.mutex.Lock()
		defer h.mutex.Unlock()
	}

	h.Percentage = percentage
}

func (h *HaveStatePercentage) ChangeState(newState State, noLock ...bool) {
	if len(noLock) == 0 || !noLock[0] {
		h.mutex.Lock()
		defer h.mutex.Unlock()
	}

	switch newState {
	case TODO:
		h.AddPercentageValue(0, true)
	case Error, Done:
		h.AddPercentageValue(100, true)
	}

	h.State = newState
}

const (
	TODO       State = "TODO"
	InProgress State = "InProgress"
	Done       State = "Done"
	Error      State = "Error"
)
