package gfsm

type Fact struct {
	Fsm      FSM
	Src, Dst *State
	args     []interface{}
	Event    *Event
}

const (
	Before bool = true
	After  bool = false
	Enter  bool = true
	Leave  bool = false
)

type cbPos struct {
	pos, once bool
}

type callback func(Fact)
