package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"sort"
)

// Actor Node --------------------------------------------

type Node struct {
	maxElems int
	inner    *Inner
	leaf     *Leaf
}

type Inner struct {
	left    *actor.PID
	right   *actor.PID
	maxLeft int
}

type Leaf struct {
	values map[int]string
}

// Messages -----------------------------------------------

type Insert struct {
	key    int
	value  string
	caller *actor.PID
}

type Error struct {
	originalMsg interface{}
}

type Success struct {
	key         int
	value       string
	originalMsg interface{}
}

// Actions ------------------------------------------------

func (state *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Insert:
		state.insert(msg, context)
	}
}

func (state *Node) insert(msg *Insert, context actor.Context) {
	if state.inner != nil {
		switch {
		case msg.key > state.inner.maxLeft:
			context.Send(state.inner.right, msg)
		case msg.key < state.inner.maxLeft:
			context.Send(state.inner.left, msg)
		case msg.key == state.inner.maxLeft:
			context.Send(msg.caller, &Error{originalMsg: msg})
		}
	} else if state.leaf != nil {
		state.leaf.values[msg.key] = msg.value
		if len(state.leaf.values) > state.maxElems {
			// Leaf becomes inner node
			state.inner = &Inner{}
			state.inner.left = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &Node{
					maxElems: state.maxElems,
					inner:    nil,
					leaf:     &Leaf{values: make(map[int]string, state.maxElems)},
				}
			}))
			state.inner.right = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &Node{
					maxElems: state.maxElems,
					inner:    nil,
					leaf:     &Leaf{values: make(map[int]string, state.maxElems)},
				}
			}))
			keys := make([]int, state.maxElems)
			for k := range state.leaf.values {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			indexMaxLeft := state.maxElems / 2
			state.inner.maxLeft = keys[indexMaxLeft]
			for _, k := range keys {
				if k <= indexMaxLeft {

				} else {

				}
			}

		}
	}
}
