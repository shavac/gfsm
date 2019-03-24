package gfsm

import "github.com/pkg/errors"

type Event struct {
	name string
	srcs []string
	dst  string
}

func (e Event) String() string {
	return e.name
}

type eFactory struct {
	e *Event
}

func EventFactory(name string) *eFactory {
	return &eFactory{e: &Event{
		name: name,
		srcs: []string{},
	}}
}

func (ef *eFactory) Build() (*Event, error) {
	if len(ef.e.srcs) == 0 || ef.e.dst == "" {
		return nil, errors.New("Cant build event")
	}
	return ef.e, nil
}

func (ef *eFactory) AddSourceStates(sts ...string) *eFactory {
	for _, st := range sts {
		if len(st) == 0 {
			continue
		}
		ef.e.srcs = append(ef.e.srcs, st)
	}
	return ef
}

func (ef *eFactory) SetDestState(dst string) *eFactory {
	if len(dst) == 0 {
		return ef
	}
	ef.e.dst = dst
	return ef
}
