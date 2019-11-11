package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"math"
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

// Actions ------------------------------------------------

func (state *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Insert:
		state.insert(msg, context)
	case *messages.Search:
		state.search(msg, context)
	case *messages.Delete:
		state.delete(msg, context)
	case *messages.UpdateMaxleft:
		state.inner.maxLeft = msg.NewValue
	case messages.Travers:
		state.travers(&msg, context)

	}
}

func (state *Node) insert(msg *messages.Insert, context actor.Context) {
	if state.inner != nil {
		switch {
		case msg.Key > state.inner.maxLeft:
			context.RequestWithCustomSender(state.inner.right, msg, context.Sender())
		case msg.Key < state.inner.maxLeft:
			context.RequestWithCustomSender(state.inner.left, msg, context.Sender())
		case msg.Key == state.inner.maxLeft:
			context.Send(context.Sender(), &messages.Error{OriginalMsg: msg})
		}
	} else if state.leaf != nil {
		_, ok := state.leaf.values[msg.Key]
		if ok {
			context.Send(context.Sender(), &messages.Error{OriginalMsg: msg})
		} else {
			state.leaf.values[msg.Key] = msg.Value
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
						context.Request(state.inner.left, &messages.Insert{Key: keys[k], Value: state.leaf.values[keys[k]]})
					} else {
						context.Request(state.inner.right, &messages.Insert{Key: keys[k], Value: state.leaf.values[keys[k]]})
					}
				}
			}
			context.Send(context.Sender(), &messages.Success{OriginalMsg: msg})
		}
	}
}

func (state *Node) search(msg *messages.Search, context actor.Context) {
	if state.inner != nil { //IF is inner node
		if msg.Key > state.inner.maxLeft { // bigger -> keep searching on the right
			context.Send(state.inner.right, msg)
		} else {
			context.Send(state.inner.left, msg) // smaller -> keep searching on the left
		}
	} else { // IF leaf
		elem, ok := state.leaf.values[msg.Key]
		if ok {
			context.Send(msg.Caller, messages.Success{
				Key:         msg.Key,
				Value:       elem,
				OriginalMsg: msg,
			})
		} else { //Key not in Tree
			context.Send(msg.Caller, messages.Error{OriginalMsg: msg})
		}
	}
}

func (node *Node) delete(msg *messages.Delete, context actor.Context) {
	if node.inner != nil { //IF is inner node
		if msg.Key <= node.inner.maxLeft { // search on left
			if msg.Key == node.inner.maxLeft { // update Maxleft
				msg.NeedUpdate = append(msg.NeedUpdate, context.Self())
			}
			context.Send(node.inner.left, msg)
		} else { // search on right
			context.Send(node.inner.right, msg)
		}
	} else { //IF is leaf
		_, OK := node.leaf.values[msg.Key]
		if OK {
			delete(node.leaf.values, msg.Key)
			maxval := 0
			for val := range node.leaf.values {
				maxval = int(math.Max(float64(val), float64(maxval)))
			}
			for _, node := range msg.NeedUpdate {
				context.Send(node, messages.UpdateMaxleft{NewValue: 1}) //TODO update with real value
			}
		} else {
			context.Send(msg.Caller, messages.Error{OriginalMsg: msg})
		}
	}
}

func (node *Node) travers(msg *messages.Travers, context actor.Context) {
	if node.inner != nil { //IF is inner node
		context.Send(node.inner.right, messages.Travers{
			Caller:     msg.Caller,
			TreeValues: nil,
		})
		context.Send(node.inner.left, messages.Travers{
			Caller:     msg.Caller,
			TreeValues: nil,
		})
	} else { //IF is leaf
		context.Send(msg.Caller, messages.Travers{
			Caller:     msg.Caller,
			TreeValues: node.leaf.values,
		})
	}
}
