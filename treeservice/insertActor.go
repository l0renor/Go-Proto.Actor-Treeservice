package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/tree"
)

type InsertActor struct {
	CLI *actor.PID
}

func (actor *InsertActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case tree.Success:
		context.Send(actor.CLI, messages.Search_Response{
			Success: true,
			Value:   msg.Value,
			Error:   "",
		}) //TODO
	case tree.Error:
		context.Send(actor.CLI, messages.Search_Response{
			Success: false,
			Value:   "",
			Error:   "",
		})
	}

}
