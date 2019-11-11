package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/tree"
)

type Service struct {
	trees map[int32]Tree
}

type Tree struct {
	Root  *actor.PID
	Token string
}

func (service *Service) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case tree.Success:
		handleSuccess(msg)
	case tree.Error:
		handleError(msg)
	}
}
