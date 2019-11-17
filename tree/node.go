package tree

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/logger"
	"sort"
)

// Actor Node --------------------------------------------

type Node struct {
	maxElems int32
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

// Constructor

func NewRoot(maxElems int32) *Node {
	return &Node{
		maxElems: maxElems,
		inner:    nil,
		leaf:     &Leaf{values: make(map[int32]string, maxElems)},
	}
}

// Actions ------------------------------------------------

func (node *Node) Receive(context actor.Context) {
	logger.GetInstance().Info.Printf("NODE %v got msg %T", node, context.Message())
	switch msg := context.Message().(type) {
	case *Insert:
		node.insert(msg, context)
	case *Search:
		node.search(msg, context)
	case *Delete:
		node.delete(msg, context)
	case *UpdateMaxLeft:
		node.inner.maxLeft = msg.NewValue
	case *Traverse:
		node.traverse(msg, context)
	case *Kill:
		node.kill(context)
	}
}

func (node *Node) insert(msg *Insert, context actor.Context) {
	logger.GetInstance().Info.Printf("Insert %v, on %v\n", msg.Key, node)
	if node.inner != nil {
		switch {
		case msg.Key > node.inner.maxLeft:
			context.RequestWithCustomSender(node.inner.right, msg, context.Sender())
		case msg.Key < node.inner.maxLeft:
			context.RequestWithCustomSender(node.inner.left, msg, context.Sender())
		case msg.Key == node.inner.maxLeft:
			context.Respond(Error{})
			logger.GetInstance().Error.Printf("key == maxleft Insert %v, on %v\n", msg, node)
		}
	} else if node.leaf != nil {
		_, ok := node.leaf.values[msg.Key]
		if ok {
			context.Respond(&Error{})
		} else {
			node.leaf.values[msg.Key] = msg.Value
			if int32(len(node.leaf.values)) > node.maxElems {
				logger.GetInstance().Info.Println("Leaf has more than max elems, new nodes have to be created")
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
				keys := make([]int, 0)
				for k := range node.leaf.values {
					keys = append(keys, int(k))
				}
				sort.Ints(keys)
				indexMaxLeft := (node.maxElems + 1) / 2
				node.inner.maxLeft = int32(keys[indexMaxLeft])
				for i := range keys {
					var child *actor.PID
					if int32(i) <= indexMaxLeft {
						child = node.inner.left
					} else {
						child = node.inner.right
					}
					msg := &Insert{Key: int32(keys[i]), Value: node.leaf.values[int32(keys[i])]}
					logger.GetInstance().Info.Printf("Sent Insert Request %v to new child %v", msg, child)
					context.Request(child, msg)
				}
			}
			logger.GetInstance().Info.Printf("Insert successful  %v, on %v\n", msg, node)
			context.Respond(&Success{})
		}
	}
}

func (node *Node) search(msg *Search, context actor.Context) {
	logger.GetInstance().Info.Printf("search %v, on %v\n", msg, node)
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
			logger.GetInstance().Info.Printf("search found %v, on %v\n", msg, node)
			context.Respond(Success{
				Key:   msg.Key,
				Value: elem,
			})
		} else { //Key not in Tree
			logger.GetInstance().Error.Printf("search key not in tree  %v, on %v\n", msg, node)
			context.Respond(Error{})
		}
	}
}

func (node *Node) delete(msg *Delete, context actor.Context) {
	logger.GetInstance().Info.Printf("Delete %v, on %v\n", msg, node)
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
		val, OK := node.leaf.values[msg.Key]
		if OK {
			delete(node.leaf.values, msg.Key)
			maxLeft := int32(0)
			for v := range node.leaf.values {
				maxLeft = max(maxLeft, v)
			}
			context.Send(context.Parent(), UpdateMaxLeft{NewValue: maxLeft})
			logger.GetInstance().Info.Printf("Delete success %v, on %v\n", msg, node)
			context.Respond(Success{
				Key:   msg.Key,
				Value: val,
			})
		} else {
			logger.GetInstance().Error.Printf("delete ke not present %v, on %v\n", msg, node)
			context.Respond(Error{})
		}
	}
}

func (node *Node) traverse(msg *Traverse, context actor.Context) {
	logger.GetInstance().Info.Printf("Traverse %v, on %v\n", msg, node)
	if node.inner != nil { //IF is inner node
		context.RequestWithCustomSender(node.inner.right, Traverse{
			TreeValues: nil,
		}, context.Sender())
		context.RequestWithCustomSender(node.inner.left, Traverse{
			TreeValues: nil,
		}, context.Sender())
		context.Respond(TraverseWaitOneMore{})
	} else { //IF is leaf
		logger.GetInstance().Info.Printf("Traverse finshed send back to helper %v, on %v\n", msg, node)
		context.Send(context.Sender(), Traverse{
			TreeValues: node.leaf.values,
		})
	}
}

func (node *Node) kill(context actor.Context) {
	logger.GetInstance().Info.Printf("Kill on %v\n", node)
	if node.inner != nil { //IF is inner node
		context.Send(node.inner.right, Kill{})
		context.Send(node.inner.left, Kill{})
		context.Stop(context.Self())

	} else { //IF is leaf
		context.Stop(context.Self())
	}
}

func (node *Node) String() string {
	if node.inner != nil {
		return fmt.Sprintf("Inner Node\n left:%v\nright: %v \nmaxleft %v", node.inner.left, node.inner.right, node.inner.maxLeft)
	} else {
		return fmt.Sprintf("Leaf:\n"+
			"Values: %v", node.leaf.values)
	}
}

func max(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
