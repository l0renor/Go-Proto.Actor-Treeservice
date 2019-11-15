package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/tree"
	"sort"
)

type inserter struct {
	cli *actor.PID
	msg messages.Insert
}

type searcher struct {
	cli *actor.PID
	msg messages.Search
}

type deleter struct {
	cli *actor.PID
	msg messages.Delete
}

type traverser struct {
	cli           *actor.PID
	msg           messages.Traverse
	nMessagesWait int
	treemap       map[int32]string
}

func (state *inserter) Receive(context actor.Context) {
	switch _ := context.Message().(type) {
	case tree.Success:
		state.msg.Response = &messages.Insert_Response{
			Success: true,
		}
		context.Send(state.cli, state.msg)
		context.Stop(context.Self())
	case tree.Error:
		state.msg.Response = &messages.Insert_Response{
			Success: false,
			Error:   "Key already present",
		}
		context.Send(state.cli, state.msg)
		context.Stop(context.Self())
	}
}

func (state *searcher) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case tree.Success:
		state.msg.Response = &messages.Search_Response{
			Success: true,
			Value:   msg.Value,
		}
		context.Send(state.cli, state.msg)
		context.Stop(context.Self())
	case tree.Error:
		state.msg.Response = &messages.Search_Response{
			Success: false,
			Error:   "Key not found",
		}
		context.Send(state.cli, state.msg)
		context.Stop(context.Self())
	}
}

func (state *deleter) Receive(context actor.Context) {
	switch _ := context.Message().(type) {
	case tree.Success:
		state.msg.Response = &messages.Delete_Response{
			Success: true,
		}
		context.Send(state.cli, state.msg)
		context.Stop(context.Self())
	case tree.Error:
		state.msg.Response = &messages.Delete_Response{
			Success: false,
			Error:   "Key not found",
		}
		context.Send(state.cli, state.msg)
		context.Stop(context.Self())
	}
}

func (state *traverser) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case tree.Traverse:
		for k, v := range msg.TreeValues { // merge maps
			state.treemap[k] = v
		}
		state.nMessagesWait--
		if state.nMessagesWait == 0 { //all leaves have answered actor can anwser to cli and die
			treeTuple := make([]*messages.Traverse_Response_Tuple, 0)
			//sort map and make it a slice
			keys := make([]int, len(msg.TreeValues))
			for k := range msg.TreeValues {
				keys = append(keys, int(k))
			}
			sort.Ints(keys)
			for key := range keys {
				treeTuple = append(treeTuple, &messages.Traverse_Response_Tuple{
					Key:   int32(key),
					Value: msg.TreeValues[int32(key)],
				})
			}
			state.msg.Response = &messages.Traverse_Response{
				Success: true,
				Tuples:  treeTuple,
				Error:   "",
			}
			context.Send(state.cli, msg)
			context.Stop(context.Self())
		}
	case tree.TraverseWaitOneMore:
		state.nMessagesWait++
	}
}
