package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"sort"
)

type Traversaktor struct {
	cli           *actor.PID
	nMessagesWait int
	treeTupel     []tupel
}

type tupel struct {
	Key   int32
	Value string
}

func (state *Traversaktor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Travers:
		keys := make([]int, len(msg.TreeValues))
		for k := range msg.TreeValues {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)
		for key := range keys {
			state.treeTupel = append(state.treeTupel, tupel{
				Key:   int32(key),
				Value: msg.TreeValues[int32(key)],
			})
		}
		state.nMessagesWait--
		if state.nMessagesWait == 0 {
			//TODO message and cli und suezied
		}

	case TraversWaitOneMore:
		state.nMessagesWait++
	}
}
