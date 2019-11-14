package cli

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"sync"
)

type receiver struct {
	wg *sync.WaitGroup
}

func (state *receiver) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Create:
		state.create(msg)
	case *messages.Insert:
		state.insert(msg)
	case *messages.Search:
		state.search(msg)
	case *messages.Delete:
		state.delete(msg)
	case *messages.Traverse:
		state.traverse(msg)
	case *messages.Remove:
		state.remove(msg)
	}
	state.wg.Done()
}

func (state *receiver) create(msg *messages.Create) {
	if msg.Response.Success {
		fmt.Printf("ID: %d, Token: %s\n", msg.Response.Id, msg.Response.Token)
	} else {
		fmt.Println(msg.Response.Error)
	}
}

func (state *receiver) insert(msg *messages.Insert) {
	if msg.Response.Success {
		fmt.Println("Insertion successful")
	} else {
		fmt.Println(msg.Response.Error)
	}
}

func (state *receiver) search(msg *messages.Search) {
	if msg.Response.Success {
		fmt.Printf("Value: %s\n", msg.Response.Value)
	} else {
		fmt.Println(msg.Response.Error)
	}
}

func (state *receiver) delete(msg *messages.Delete) {
	if msg.Response.Success {
		fmt.Println("Deletion successful")
	} else {
		fmt.Println(msg.Response.Error)
	}
}

func (state *receiver) traverse(msg *messages.Traverse) {
	if msg.Response.Success {
		for _, tuple := range msg.Response.Tuples {
			fmt.Printf("Key: %d, Value: %s\n", tuple.Key, tuple.Value)
		}
	} else {
		fmt.Println(msg.Response.Error)
	}
}

func (state *receiver) remove(msg *messages.Remove) {
	if msg.Response.Success {
		fmt.Println("Removal successful")
	} else {
		fmt.Println(msg.Response.Error)
	}
}
