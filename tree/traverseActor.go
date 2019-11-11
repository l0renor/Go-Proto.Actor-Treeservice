package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"sort"
)

//todo service needs to send msg to root node with custom sender = traversactor
type Traversaktor struct {
	service       *actor.PID
	nMessagesWait int
	treemap       map[int32]string
}

type tuple struct {
	Key   int32
	Value string
}

func (state *Traversaktor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Travers:
		for k, v := range msg.TreeValues { // merge maps
			state.treemap[k] = v
		}
		state.nMessagesWait--
		if state.nMessagesWait == 0 { //all leaves have answered actor can anwser to service and die
			treeTuple := make([]tuple, 0)

			//sort map and make it a slice
			keys := make([]int, len(msg.TreeValues))
			for k := range msg.TreeValues {
				keys = append(keys, int(k))
			}
			sort.Ints(keys)
			for key := range keys {
				treeTuple = append(treeTuple, tuple{
					Key:   int32(key),
					Value: msg.TreeValues[int32(key)],
				})
			}
			context.Send(state.service, TraverseActor_Msg{tree: treeTuple})
			context.Stop(context.Self())
		}

	case TraversWaitOneMore:
		state.nMessagesWait++
	}
}
