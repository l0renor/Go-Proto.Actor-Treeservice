package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"sort"
)

//TODO set caller in CLI/CLI
// Actor Node --------------------------------------------

type Node struct {
	maxElems int
	inner    *Inner
	leaf     *Leaf
}

type Inner struct {
	left    *actor.PID
	right   *actor.PID
	maxLeft int32
}

type Leaf struct {
	values map[int32]string
}

// Actions ------------------------------------------------

func (node *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Insert:
		node.insert(msg, context)
	case *Search:
		node.search(msg, context)
	case *Delete:
		node.delete(msg, context)
	case *UpdateMaxLeft:
		node.inner.maxLeft = msg.NewValue
	case Travers:
		node.travers(&msg, context)

	}
}

func (node *Node) insert(msg *Insert, context actor.Context) {
	if node.inner != nil {
		switch {
		case msg.Key > node.inner.maxLeft:
			context.RequestWithCustomSender(node.inner.right, msg, context.Sender())
		case msg.Key < node.inner.maxLeft:
			context.RequestWithCustomSender(node.inner.left, msg, context.Sender())
		case msg.Key == node.inner.maxLeft:
			context.Respond(Error{OriginalMsg: msg})
		}
	} else if node.leaf != nil {
		_, ok := node.leaf.values[msg.Key]
		if ok {
			context.Respond(&Error{OriginalMsg: msg})
		} else {
			node.leaf.values[msg.Key] = msg.Value
			if len(node.leaf.values) > node.maxElems {
				// Leaf becomes inner node
				node.inner = &Inner{}
				node.inner.left = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
					return &Node{
						maxElems: node.maxElems,
						inner:    nil,
						leaf:     &Leaf{values: make(map[int32]string, node.maxElems)},
					}
				}))
				node.inner.right = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
					return &Node{
						maxElems: node.maxElems,
						inner:    nil,
						leaf:     &Leaf{values: make(map[int32]string, node.maxElems)},
					}
				}))
				keys := make([]int, node.maxElems+1)
				for k := range node.leaf.values {
					keys = append(keys, int(k))
				}
				sort.Ints(keys)
				indexMaxLeft := (node.maxElems + 1) / 2
				node.inner.maxLeft = int32(keys[indexMaxLeft])
				for _, k := range keys {
					var child *actor.PID
					if k <= indexMaxLeft {
						child = node.inner.left
					} else {
						child = node.inner.right
					}
					context.Request(child, &Insert{Key: int32(keys[k]), Value: node.leaf.values[int32(keys[k])]})
				}
			}
			context.Respond(&Success{OriginalMsg: msg})
		}
	}
}

func (node *Node) search(msg *Search, context actor.Context) {
	if node.inner != nil { //IF is inner node
		var child *actor.PID
		if msg.Key > node.inner.maxLeft { // bigger -> keep searching on the right
			child = node.inner.right
		} else {
			child = node.inner.left // smaller -> keep searching on the left
		}
		context.RequestWithCustomSender(child, msg, context.Sender())
	} else { // IF leaf
		elem, ok := node.leaf.values[msg.Key]
		if ok {
			context.Respond(Success{
				Key:         msg.Key,
				Value:       elem,
				OriginalMsg: msg,
			})
		} else { //Key not in Tree
			context.Respond(Error{OriginalMsg: msg})
		}
	}
}

func (node *Node) delete(msg *Delete, context actor.Context) {
	if node.inner != nil { //IF is inner node
		var child *actor.PID
		switch {
		case msg.Key <= node.inner.maxLeft:
			child = node.inner.left
		case msg.Key > node.inner.maxLeft:
			child = node.inner.right
		}
		context.RequestWithCustomSender(child, msg, context.Sender())
	} else if node.leaf != nil { //IF is leaf
		_, OK := node.leaf.values[msg.Key]
		if OK {
			delete(node.leaf.values, msg.Key)
			maxLeft := int32(0)
			for v := range node.leaf.values {
				maxLeft = max(maxLeft, v)
			}
			context.Send(context.Parent(), UpdateMaxLeft{NewValue: maxLeft})
		} else {
			context.Respond(Error{OriginalMsg: msg})
		}
	}
}

func (node *Node) travers(msg *Travers, context actor.Context) {
	if node.inner != nil { //IF is inner node
		context.RequestWithCustomSender(node.inner.right, Travers{
			TreeValues: nil,
		}, context.Sender())
		context.RequestWithCustomSender(node.inner.left, Travers{
			TreeValues: nil,
		}, context.Sender())
		context.Respond(TraversWaitOneMore{})
	} else { //IF is leaf
		context.Send(context.Sender(), Travers{
			TreeValues: node.leaf.values,
		})
	}
}

func max(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
