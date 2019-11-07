package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"sort"
)

//TODO set caller in CLI/service
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

type Search struct {
	caller *actor.PID
	key    int
}

type Insert struct {
	key   int
	value string
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
	case *Search:
		state.search(msg, context)
	}
}

func (state *Node) insert(msg *Insert, context actor.Context) {
	if state.inner != nil {
		switch {
		case msg.key > state.inner.maxLeft:
			context.RequestWithCustomSender(state.inner.right, msg, context.Sender())
		case msg.key < state.inner.maxLeft:
			context.RequestWithCustomSender(state.inner.left, msg, context.Sender())
		case msg.key == state.inner.maxLeft:
			context.Send(context.Sender(), &Error{originalMsg: msg})
		}
	} else if state.leaf != nil {
		_, ok := state.leaf.values[msg.key]
		if ok {
			context.Send(context.Sender(), &Error{originalMsg: msg})
		} else {
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
				keys := make([]int, state.maxElems+1)
				for k := range state.leaf.values {
					keys = append(keys, k)
				}
				sort.Ints(keys)
				indexMaxLeft := (state.maxElems + 1) / 2
				state.inner.maxLeft = keys[indexMaxLeft]
				for _, k := range keys {
					if k <= indexMaxLeft {
						context.Request(state.inner.left, &Insert{key: keys[k], value: state.leaf.values[keys[k]]})
					} else {
						context.Request(state.inner.right, &Insert{key: keys[k], value: state.leaf.values[keys[k]]})
					}
				}
			}
			context.Send(context.Sender(), &Success{originalMsg: msg})
		}
	}
}

func (state *Node) search(msg *Search, context actor.Context) {
	if state.inner != nil { //IF is inner node
		if msg.key > state.inner.maxLeft { // bigger -> keep searching on the right
			context.Send(state.inner.right, msg)
		} else {
			context.Send(state.inner.left, msg) // smaller -> keep searching on the left
		}
	} else { // IF leaf
		elem, ok := state.leaf.values[msg.key]
		if ok {
			context.Send(msg.caller, Success{
				key:         msg.key,
				value:       elem,
				originalMsg: msg,
			})
		} else { //Key not in Tree
			context.Send(msg.caller, Error{originalMsg: msg})
		}
	}
}
