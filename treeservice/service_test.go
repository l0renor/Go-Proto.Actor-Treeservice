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
var token2 string

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
		token2 = msg.Response.Token
	}
	globl_msg = nil
}

func TestInsert_Traverse(t *testing.T) {
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
	for i := 2; i < 10; i++ { //insert 2 to 9 again should result in Key already present error
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
			if msg.Success {
				t.Error("Insert -> Response -> sucsess should bee false")
			}
			if msg.Error != "Key already present" {
				t.Error("Insert wrong error got " + msg.Error)
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
		want := [8]*messages.Traverse_Response_Tuple{{Key: 2, Value: "2"}, &messages.Traverse_Response_Tuple{Key: 3, Value: "3"}, &messages.Traverse_Response_Tuple{Key: 4, Value: "4"}, &messages.Traverse_Response_Tuple{Key: 5, Value: "5"}, &messages.Traverse_Response_Tuple{Key: 6, Value: "6"}, &messages.Traverse_Response_Tuple{Key: 7, Value: "7"}, &messages.Traverse_Response_Tuple{Key: 8, Value: "8"}, &messages.Traverse_Response_Tuple{Key: 9, Value: "9"}}
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
	//search on empty
	for i := 10; i < 20; i++ {
		f := context.RequestFuture(servicePID, &messages.Search{
			Id:       2,
			Token:    token2,
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
func TestDelete(t *testing.T) {
	for i := 2; i < 10; i++ { //delete present key -> sucsess
		f := context.RequestFuture(servicePID, &messages.Delete{
			Id:       1,
			Token:    token,
			Key:      int32(i),
			Response: nil,
		}, 500*time.Millisecond)

		res, err := f.Result()
		if err != nil {
			t.Error("Timeout delete Tree")
		}

		switch msg := res.(type) {
		case messages.Delete:
			if !msg.Response.Success {
				t.Error(" Delete not successful Key was present")
			}
		}
	}
	for i := 2; i < 10; i++ { //delete missing key -> error
		f := context.RequestFuture(servicePID, &messages.Delete{
			Id:       1,
			Token:    token,
			Key:      int32(i),
			Response: nil,
		}, 500*time.Millisecond)

		res, err := f.Result()
		if err != nil {
			t.Error("Timeout delete Tree")
		}

		switch msg := res.(type) {
		case messages.Delete:
			if msg.Response.Success {
				t.Error(" Delete  successful Key was NOT present")
			}
		}
	}

}

func TestRemove_Creds(t *testing.T) {
	//remove present Tree
	f := context.RequestFuture(servicePID, &messages.Remove{
		Id:       1,
		Token:    token,
		Response: nil,
	}, 500*time.Millisecond)

	res, err := f.Result()
	if err != nil {
		t.Error("Timeout remove Tree")
	}

	switch msg := res.(type) {
	case messages.Remove:
		if !msg.Response.Success {
			t.Error(" Remove not  successful")
		}
	}
	//remove not present tree
	f = context.RequestFuture(servicePID, &messages.Remove{
		Id:       1,
		Token:    token,
		Response: nil,
	}, 500*time.Millisecond)

	res, err = f.Result()
	if err != nil {
		t.Error("Timeout remove Tree")
	}

	switch msg := res.(type) {
	case messages.Remove:
		if msg.Response.Success {
			t.Error(" Remove successful  Tree was not present")
		}
	}

	//remove present  with wrong creds
	f = context.RequestFuture(servicePID, &messages.Remove{
		Id:       2,
		Token:    token,
		Response: nil,
	}, 500*time.Millisecond)

	res, err = f.Result()
	if err != nil {
		t.Error("Timeout remove Tree")
	}

	switch msg := res.(type) {
	case messages.Remove:
		if msg.Response.Success {
			t.Error(" Remove successful  wrong creds")
		}
	}
}

func TestBigInsert_delete(t *testing.T) {
	f := context.RequestFuture(servicePID, &messages.Create{
		MaxElems: 3,
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
		if msg.Response.Token == "" {
			t.Error("Token missing")
		}
		token = msg.Response.Token
	}
	//Insert round  1
	round1 := []int32{7, 9, 15, 20, 22, 30}

	for i := range round1 {
		f := context.RequestFuture(servicePID, &messages.Insert{
			Id:       3,
			Token:    token,
			Key:      round1[i],
			Value:    strconv.Itoa(int(round1[i])),
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
	//delete some for structurl change
	del1 := []int32{9, 15, 22, 30}

	for i := range del1 {
		f := context.RequestFuture(servicePID, &messages.Delete{
			Id:       1,
			Token:    token,
			Key:      del1[i],
			Response: nil,
		}, 500*time.Millisecond)

		res, err := f.Result()
		if err != nil {
			t.Error("Timeout delete Tree")
		}

		switch msg := res.(type) {
		case messages.Delete:
			if !msg.Response.Success {
				t.Error(" Delete not successful Key was present")
			}
		}
	}

	round2 := []int32{2, 40, 21, 35}

	for i := range round2 {
		f := context.RequestFuture(servicePID, &messages.Insert{
			Id:       3,
			Token:    token,
			Key:      round2[i],
			Value:    strconv.Itoa(int(round2[i])),
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

	f = context.RequestFuture(servicePID, &messages.Traverse{
		Id:       1,
		Token:    token,
		Response: nil,
	}, 500*time.Millisecond)

	res, err = f.Result()
	if err != nil {
		t.Error("Timeout traverse Tree")
	}

	switch msg := res.(type) {
	case messages.Traverse:
		if !msg.Response.Success {
			t.Error("Traverse not successful " + msg.Response.Error)
		}
		want := [6]*messages.Traverse_Response_Tuple{{Key: 2, Value: "2"}, &messages.Traverse_Response_Tuple{Key: 7, Value: "7"}, &messages.Traverse_Response_Tuple{Key: 20, Value: "20"}, &messages.Traverse_Response_Tuple{Key: 21, Value: "21"}, &messages.Traverse_Response_Tuple{Key: 35, Value: "35"}, &messages.Traverse_Response_Tuple{Key: 40, Value: "40"}}
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
