package service

import (
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
			t.Error("Create -> Response -> Success == false")
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
	for i := 2; i < 10; i++ { //insert 2 to 9
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
		want := [8]*messages.Traverse_Response_Tuple{&messages.Traverse_Response_Tuple{Key: 2, Value: "2"}, &messages.Traverse_Response_Tuple{Key: 3, Value: "3"}, &messages.Traverse_Response_Tuple{Key: 4, Value: "4"}, &messages.Traverse_Response_Tuple{Key: 5, Value: "5"}, &messages.Traverse_Response_Tuple{Key: 6, Value: "6"}, &messages.Traverse_Response_Tuple{Key: 7, Value: "7"}, &messages.Traverse_Response_Tuple{Key: 8, Value: "8"}, &messages.Traverse_Response_Tuple{Key: 9, Value: "9"}}
		for e := range want {
			if want[e].Key != msg.Response.Tuples[e].Key {
				t.Errorf("Missmatch in traverse Key w:%v, is: %v", want[e].Key, msg.Response.Tuples[e].Key)
			}
			if want[e].Value != msg.Response.Tuples[e].Value {
				t.Errorf("Missmatch in traverse Value w:%v, is: %v", want[e].Value, msg.Response.Tuples[e].Value)
			}
		}
	}
}

func TestSearch(t *testing.T) {
	// search present keys
	for i := 2; i < 10; i++ {
		f := context.RequestFuture(servicePID, &messages.Search{
			Id:       1,
			Token:    token,
			Key:      int32(i),
			Response: nil,
		}, 500*time.Millisecond)

		res, err := f.Result()
		if err != nil {
			t.Error("Timeout search Tree")
		}

		switch msg := res.(type) {
		case messages.Search:
			if !msg.Response.Success {
				t.Error("Search Response not successful Key was present")
			}
			if msg.Response.Value != strconv.Itoa(i) {
				t.Errorf("Wrong value in search want %v got %v", strconv.Itoa(i), msg.Response.Value)
			}
		}
	}
	//search non present keys
	for i := 10; i < 20; i++ {
		f := context.RequestFuture(servicePID, &messages.Search{
			Id:       1,
			Token:    token,
			Key:      int32(i),
			Response: nil,
		}, 500*time.Millisecond)

		res, err := f.Result()
		if err != nil {
			t.Error("Timeout search Tree")
		}

		switch msg := res.(type) {
		case messages.Search:
			if msg.Response.Success {
				t.Error("Search Response  successful Key was NOT present")
			}
			if msg.Response.Error != "Key not found" {
				t.Error("Wrong error n key not found: " + msg.Response.Error)
			}
		}
	}

}
