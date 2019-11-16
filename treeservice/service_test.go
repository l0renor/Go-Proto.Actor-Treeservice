package service

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"testing"
	"time"
)

type tester struct {
}

var globl_msg interface{}
var servicePID *actor.PID
var testerPID *actor.PID
var context *actor.RootContext
var token string

func (test *tester) Receive(context actor.Context) {
	globl_msg = context.Message()
}
func TestInitService(t *testing.T) {
	context = actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Service{}
	})
	servicePID = context.Spawn(props)
}

func TestCreate(t *testing.T) {
	f := context.RequestFuture(servicePID, &messages.Create{
		MaxElems: 2,
		Response: nil,
	}, 100*time.Millisecond)

	res, err := f.Result()
	if err != nil {
		t.Error("Timeout create Tree")
	}
	switch msg := res.(type) {
	case *messages.Create:
		if !msg.Response.Success {
			t.Error("Create -> Response -> sucsess == false")
		}
		if msg.Response.Id != 1 {
			t.Errorf("Initial id  should be 1 was %v", msg.Response.Id)
		}
		if msg.Response.Token == "" {
			t.Error("Token missing")
		}
		token = msg.Response.Token
	}
	//2. test
	f = context.RequestFuture(servicePID, &messages.Create{
		MaxElems: 2,
		Response: nil,
	}, 100*time.Millisecond)

	res, err = f.Result()
	if err != nil {
		t.Error("Timeout create Tree")
	}
	switch msg := res.(type) {
	case *messages.Create:
		if !msg.Response.Success {
			t.Error("Create -> Response -> sucsess == false")
		}
		if msg.Response.Id != 2 {
			t.Errorf("Second id  should be 2 was %v", msg.Response.Id)
		}
		if msg.Response.Token == "" {
			t.Error("Token missing")
		}
	}
}
func TestInsert(t *testing.T) {
	for i := 0; i < 1; i++ { //insert 0 to 9
		context.Send(servicePID, messages.Insert{
			Id:       1,
			Token:    token,
			Key:      int32(i),
			Value:    string(i),
			Response: nil,
		})
		time.Sleep(time.Millisecond * 100)
		switch msg := globl_msg.(type) {
		case *messages.Insert:
			if !msg.Response.Success {
				t.Error("Insert -> Response -> sucsess == false")
			}
		}
	}
	//validate that the values are preseent
	context.Send(servicePID, &messages.Traverse{
		Id:       1,
		Token:    token,
		Response: nil,
	})
	fmt.Print(globl_msg)
	switch msg := globl_msg.(type) {

	case *messages.Traverse:
		if !msg.Response.Success {
			t.Error("Insert -> Traverse -> sucsess == false")
		}
		print(msg.Response.Tuples)
	}

}
