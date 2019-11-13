package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/tree"
	"sort"
)

type InsertActor struct {
	CLI *actor.PID
}

func (service *InsertActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case tree.Success:
		context.Send(service.CLI, INSERTSUCSESSTOCLI) //TODO
	case tree.Error:
		//send error to cli

	}

}
