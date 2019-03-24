// +build ignore

package main

import (
	"fmt"
	fsm "github.com/shavac/gfsm"
)

var (
	evtOpen, evtClose *fsm.Event
	door              fsm.FSM
)

func init() {
	evtOpen, _ = fsm.EventFactory("open").
		AddSourceStates("closed").
		SetDestState("opened").
		Build()
	evtClose, _ = fsm.EventFactory("close").
		AddSourceStates("opened").SetDestState("closed").Build()
	door, _ = fsm.FSMFactory("door").
		RegisterEvents(evtOpen, evtClose).
		RegisterGlobalCallback(fsm.Before,
			func(fct fsm.Fact) {
				fmt.Printf("gfsm %s transiting from state %s to %s\n", fct.Fsm, fct.Src, fct.Dst)
			}).
		RegisterStateCallback("opened", fsm.Enter,
			func(fct fsm.Fact) {
				fmt.Printf("gfsm %s entered opened state\n", fct.Fsm)
			}).
		Initial("closed").
		Build()
}

func main() {
	wait := make(chan bool)
	go func() {
		waitEnter := make(chan bool)
		door.RegisterStateTempCallback("opened", fsm.Enter, func(fct fsm.Fact) {
			fmt.Println("in temp callback")
			waitEnter <- true
		})
		<-waitEnter
		fmt.Println("ok...i have entered.")
		wait <- true
	}()


	fmt.Println("current state is", door.Current())
	err := door.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("current state is", door.Current())

	err = door.Event("close")
	if err != nil {
		fmt.Println(err)
	}

	err = door.Event("open")
	if err != nil {
		fmt.Println(err)
	}
	err = door.Event("close")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("current state is", door.Current())
	<-wait
}
