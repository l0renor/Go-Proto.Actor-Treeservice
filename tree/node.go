package tree

import "github.com/AsynkronIT/protoactor-go/actor"

//TODO set caller in CLI/
type Node struct {
	maxValues int
	inner     Inner
	leaf      Leaf
}

type Inner struct {
	left    actor.PID
	right   actor.PID
	maxLeft int
}

type Leaf struct {
	values map[int]string
}

type searchMSG struct {
	caller actor.PID
	key    int
}

func (node *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *searchMSG:
		node.search(msg, context)

	}
}

func (node *Node) search(msg *searchMSG, context actor.Context) {
	if (node.inner != Inner{}) { //IF is inner node
		if msg.key > node.inner.maxLeft { // bigger -> keep searching on the right
			context.Send(&node.inner.right, msg)
		} else {
			context.Send(&node.inner.left, msg) // smaller -> keep searching on the left
		}
	} else { // IF leaf
		elem, ok := node.leaf.values[msg.key]
		if ok {
			context.Send(msg.caller)
		}
	}
}
