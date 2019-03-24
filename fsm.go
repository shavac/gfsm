package gfsm

import (
	"errors"
	"fmt"
	_ "log"
	"sync"
)

type FSM interface {
	fmt.Stringer
	Current() *State
	SetState(string) error
	Event(string, ...interface{}) error
	RegisterGlobalTempCallback(bool, callback) FSM
	RegisterStateTempCallback(string, bool, callback) FSM
}

const (
	before = 0
	enter = 0
	after = 1
	leave = 1
	temp_before = 2
	temp_enter = 2
	temp_after = 3
	temp_leave = 3
)

type fsm struct {
	name         string
	current      string
	states       map[string]*State
	// 0: before/enter, 1: after/leave, 2: temp before/enter 3: temp after/leave
	callbacks    map[string][4][]callback
	stLock       *sync.RWMutex
	evtLock      *sync.Mutex
}

func (m fsm) String() string {
	return m.name
}

func (m *fsm) Current() *State {
	return m.states[m.current]
}

func (m *fsm) SetState(stName string) error {
	m.stLock.Lock()
	if _, ok := m.states[stName]; !ok {
		return errors.New("state does not exist")
	}
	m.current = stName
	m.stLock.Unlock()
	return nil
}

func (m *fsm) GetState(stName string) *State {
	if st, ok := m.states[stName]; !ok {
		return nil
	} else {
		return st
	}
}

func (m *fsm) registerCallback(stName string, posIdx int, f callback) FSM {
	cbMatrix := m.callbacks[stName]
	cbMatrix[posIdx] = append(cbMatrix[posIdx], f)
	m.callbacks[stName] = cbMatrix
	return m
}

func (m *fsm) RegisterStateTempCallback(stName string, pos bool, f callback) FSM {
	posIdx := 0
	if pos {
		posIdx = temp_enter
	} else {
		posIdx = temp_leave
	}
	return m.registerCallback(stName, posIdx, f)
}

func (m *fsm) RegisterGlobalTempCallback(pos bool, f callback) FSM {
	return m.RegisterStateTempCallback("", pos, f)
}

func (m *fsm) RegisterStateCallback(stName string, pos bool, f callback) FSM {
	posIdx := 0
	if pos {
		posIdx = 0
	} else {
		posIdx = 1
	}
	return m.registerCallback(stName, posIdx, f)
}

func (m *fsm) RegisterGlobalCallback(pos bool, f callback) FSM {
	return m.RegisterStateCallback("", pos, f)
}

func (m *fsm) Event(evtName string, args ...interface{}) error {
	e, ok := m.Current().events[evtName]
	if !ok {
		return errors.New("cannot raise event in this state")
	}
	m.evtLock.Lock()
	defer m.evtLock.Unlock()
	fct := Fact{
		Fsm:   m,
		Src:   m.Current(),
		Dst:   m.states[e.dst],
		args:  args,
		Event: e,
	}
	gcb := m.callbacks[""]
	scb := m.callbacks[m.current]
	for _, f := range gcb[temp_before] {
		f(fct)
	}
	for _, f := range gcb[before] {
		f(fct)
	}

	for _, f := range scb[temp_leave] {
		f(fct)
	}
	for _, f := range scb[leave] {
		f(fct)
	}

	gcb[temp_before] = []callback{}
	m.callbacks[""] = gcb
	scb[temp_leave] = []callback{}
	m.callbacks[m.current] = scb

	m.current = e.dst

	scb = m.callbacks[m.current]

	for _, f := range gcb[temp_after] {
		f(fct)
	}
	for _, f := range gcb[after] {
		f(fct)
	}

	for _, f := range scb[temp_enter] {
		f(fct)
	}
	for _, f := range scb[enter] {
		f(fct)
	}

	gcb[temp_after] = []callback{}
	m.callbacks[""] = gcb

	scb[temp_enter] = []callback{}
	m.callbacks[m.current] = scb
	return nil
}

type mFactory struct {
	m *fsm
}

func FSMFactory(name string) *mFactory {
	stl := &sync.RWMutex{}
	evtl := &sync.Mutex{}
	mf := &mFactory{&fsm{
		name:         name,
		states:       make(map[string]*State),
		callbacks:    make(map[string][4][]callback),
		stLock:       stl,
		evtLock:      evtl,
	}}
	mf.m.callbacks[""] = [4][]callback{}
	return mf
}

func (mf *mFactory) Initial(stName string) *mFactory {
	mf.m.current = stName
	return mf
}

func (mf *mFactory) RegisterEvents(evts ...*Event) *mFactory {
	for _, evt := range evts {
		for _, st := range evt.srcs {
			if _, ok := mf.m.states[st]; !ok {
				mf.m.states[st] = newState(st)
				mf.m.callbacks[st] = [4][]callback{}
			}
			mf.m.states[st].events[evt.name] = evt
		}
		if _, ok := mf.m.states[evt.dst]; !ok {
			mf.m.states[evt.dst] = newState(evt.dst)
		}
	}
	return mf
}

func (mf *mFactory) RegisterStateCallback(stName string, pos bool, f callback) *mFactory {
	posIdx := 0
	if pos {
		posIdx = before
	} else {
		posIdx = after
	}
	mf.m.registerCallback(stName, posIdx, f)
	return mf
}

func (mf *mFactory) RegisterGlobalCallback(pos bool, f callback) *mFactory {
	return mf.RegisterStateCallback("", pos, f)
}

func (mf *mFactory) Build() (*fsm, error) {
	if _, ok := mf.m.states[mf.m.current]; !ok {
		return nil, errors.New("initial state " + mf.m.current + " wrong")
	}
	return mf.m, nil
}
