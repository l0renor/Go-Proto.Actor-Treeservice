package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"testing"
	"time"
)

type tester struct {
}

var msg interface{}

func (tester *tester) Receive(context actor.Context) {
	msg = context.Message()
}

func TestCreate(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Service{}
	})
	servicePid := context.Spawn(props)

	props = actor.PropsFromProducer(func() actor.Actor {
		return &tester{}
	})
	testerPid := context.Spawn(props)

	context.RequestWithCustomSender(servicePid, &messages.Create{
		MaxElems: 3,
		Response: nil,
	}, testerPid)
	time.Sleep(100 * time.Millisecond)
	switch thismsg := msg.(type) {
	case *messages.Create:
		if thismsg.Response.Id != 1 {
			t.Errorf("Wrong id wanted 1 got %v", thismsg.Response.Id)
		}

	}

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
