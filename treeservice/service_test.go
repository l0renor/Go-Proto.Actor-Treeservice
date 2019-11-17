package service

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"strconv"
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
	globl_msg = nil
}

func TestInsert(t *testing.T) {
	for i := 2; i < 4; i++ { //insert 2 to 9
		f := context.RequestFuture(servicePID, &messages.Insert{
			Id:       1,
			Token:    token,
			Key:      int32(i),
			Value:    strconv.Itoa(i),
			Response: nil,
		}, 100*time.Millisecond)

		res, err := f.Result()
		if err != nil {
			t.Error("Timeout insert Tree")
		}

		switch msg := res.(type) {
		case *messages.Insert_Response:
			if !msg.Success {
				t.Error("Insert -> Response -> sucsess == false")
			}
		}
	}
	f := context.RequestFuture(servicePID, &messages.Traverse{
		Id:       1,
		Token:    token,
		Response: nil,
	}, 500*time.Millisecond)

	res, err := f.Result()
	if err != nil {
		t.Error("Timeout traverse Tree")
	}

	switch msg := res.(type) {
	case messages.Traverse:
		if !msg.Response.Success {

			t.Error("Traverse not successful " + msg.Response.Error)
		}
		fmt.Sprintf("Traverse: %v", msg.Response.Tuples)

	}
}

//func TestInsert(t *testing.T) {
//	for i := 0; i < 1; i++ { //insert 0 to 9
//		context.Send(servicePID, &messages.Insert{
//			Id:       1,
//			Token:    token,
//			Key:      int32(i),
//			Value:    string(i),
//			Response: nil,
//		})
//		time.Sleep(time.Millisecond * 500)
//		if globl_msg == nil {
//			t.Errorf("Insert Response missing %v",i)
//		}
//		switch msg := globl_msg.(type) {
//		case *messages.Insert:
//			if !msg.Response.Success {
//				t.Error("Insert -> Response -> sucsess == false")
//			}
//		}
//		globl_msg = nil
//	}
//	//validate that the values are preseent
//	context.Send(servicePID, &messages.Traverse{
//		Id:       1,
//		Token:    token,
//		Response: nil,
//	})
//	fmt.Print(globl_msg)
//	switch msg := globl_msg.(type) {
//
//	case *messages.Traverse:
//		if !msg.Response.Success {
//			t.Error("Insert -> Traverse -> sucsess == false")
//		}
//		print(msg.Response.Tuples)
//	}
//	globl_msg = nil
//}
