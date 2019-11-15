package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	context := actor.EmptyRootContext
	testService := &Service{}
	props := actor.PropsFromProducer(func() actor.Actor {
		return testService
	})
	servicePid := context.Spawn(props)
	context.Send(servicePid, messages.Create{
		MaxElems: 2,
		Response: nil,
	})
	time.Sleep(2 * time.Second)
	_, ok := testService.trees[1]
	if !ok {
		t.Error("No Tree was created")
	}

}
