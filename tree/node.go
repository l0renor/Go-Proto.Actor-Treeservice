package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

//TODO set caller in CLI/service
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
type updateMaxleft struct {
	newValue int
}
type Search struct {
	caller *actor.PID
	key    int
}

type Insert struct {
	key    int
	value  string
	caller *actor.PID
}
type Delete struct {
	key        int
	caller     *actor.PID
	needUpdate []*actor.PID
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
	case *Delete:
		state.delete(msg, context)
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

func (node *Node) search(msg *Search, context actor.Context) {
	if (node.inner != Inner{}) { //IF is inner node
		if msg.key > node.inner.maxLeft { // bigger -> keep searching on the right
			context.Send(node.inner.right, msg)
		} else {
			context.Send(node.inner.left, msg) // smaller -> keep searching on the left
		}
	} else { // IF leaf
		elem, ok := node.leaf.values[msg.key]
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
func (node *Node) delete(msg *Delete, context actor.Context) {
	if (node.inner != Inner{}) { //IF is inner node
		if msg.key <= node.inner.maxLeft { // search on left
			if msg.key == node.inner.maxLeft { // update Maxleft
				msg.needUpdate = append(msg.needUpdate, context.Self())
			}
			context.Send(node.inner.left, msg)
		} else { // search on right
			context.Send(node.inner.right, msg)
		}
	} else { //IF is leaf
		_, OK := node.leaf.values[msg.key]
		if OK {
			delete(node.leaf.values, msg.key)
			for _, node := range node.leaf.values {
				context.Send(node, updateMaxleft{newValue: 1}) //TODO update with real value
			}
		} else {
			context.Send(msg.caller, Error{originalMsg: msg})
		}
	}
}
