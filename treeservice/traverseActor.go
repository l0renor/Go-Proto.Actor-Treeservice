package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/tree"
	"sort"
)

//todo CLI needs to send msg to root node with custom sender = traversactor
type Traversaktor struct {
	CLI           *actor.PID
	NMessagesWait int
	treemap       map[int32]string
}

type Tuple struct {
	Key   int32
	Value string
}

func (state *Traversaktor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case tree.Travers:
		for k, v := range msg.TreeValues { // merge maps
			state.treemap[k] = v
		}
		state.NMessagesWait--
		if state.NMessagesWait == 0 { //all leaves have answered actor can anwser to CLI and die
			treeTuple := make([]Tuple, 0)

			//sort map and make it a slice
			keys := make([]int, len(msg.TreeValues))
			for k := range msg.TreeValues {
				keys = append(keys, int(k))
			}
			sort.Ints(keys)
			for key := range keys {
				treeTuple = append(treeTuple, Tuple{
					Key:   int32(key),
					Value: msg.TreeValues[int32(key)],
				})
			}
			//context.Send(state.CLI, tree.TraverseActor_Msg{tree: treeTuple}) TODO send protomsg to cli
			context.Stop(context.Self())
		}

	case tree.TraversWaitOneMore:
		state.NMessagesWait++
	}
}
