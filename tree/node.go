package tree

import "github.com/AsynkronIT/protoactor-go/actor"

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
