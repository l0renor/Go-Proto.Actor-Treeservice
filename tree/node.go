package tree

import "github.com/AsynkronIT/protoactor-go/actor"

// Actor Node --------------------------------------------

type Node struct {
	maxValues int
	inner     Inner
	leaf      Leaf
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
	if state.inner != (Inner{}) {
		switch {
		case msg.key > state.inner.maxLeft:
			context.Send(state.inner.right, msg)
		case msg.key < state.inner.maxLeft:
			context.Send(state.inner.left, msg)
		case msg.key == state.inner.maxLeft:
			context.Send(msg.caller, &Error{originalMsg: msg})
		}
	}
}
