package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Service{}
	})
	servicePid := context.Spawn(props)

	context.Send(servicePid, messages.Create{
		MaxElems: 2,
		Response: nil,
	})

	time.Sleep(2 * time.Second)

	//f := context.RequestFuture(servicePid,messages.Create{
	//	MaxElems: 2,
	//	Response: nil,
	//},2*time.Second)
	//
	//res, err := f.Result()
	//	if err!=nil {
	//	t.Error("Timeout create Tree")
	//}
	//switch msg := res.(type) {
	//case *messages.Create:
	//	if !msg.Response.Success{
	//		t.Error("Create -> Response -> sucsess == false")
	//	}
	//	if msg.Response.Id != 1{
	//		t.Errorf("Initial id  should be 1 was %v",msg.Response.Id)
	//	}
	//}

}
