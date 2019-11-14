package cli

import (
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
}

func (state *receiver) create(create *messages.Create) {

}

func (state *receiver) insert(insert *messages.Insert) {

}

func (state *receiver) search(search *messages.Search) {

}

func (state *receiver) delete(msg *messages.Delete) {

}

func (state *receiver) traverse(traverse *messages.Traverse) {

}

func (state *receiver) remove(remove *messages.Remove) {

}
