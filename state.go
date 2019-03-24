package gfsm

type State struct {
	name   string
	events map[string]*Event
}

func newState(name string) *State {
	return &State{
		name,
		make(map[string]*Event),
	}
}

func (st State) String() string {
	return st.name
}

func (st State) DestStates() []string {
	trs := []string{}
	for _, e := range st.events {
		trs = append(trs, e.dst)
	}
	return trs
}

func (st State) CanTransitTo(stName string) bool {
	for _, d := range st.DestStates() {
		if d == stName {
			return true
		}
	}
	return false
}
