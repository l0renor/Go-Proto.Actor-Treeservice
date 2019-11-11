package messages

import "github.com/AsynkronIT/protoactor-go/actor"

type UpdateMaxleft struct {
	NewValue int
}

//Die nachricht verläuft den baum nach unten und teilt sich auf so oft nötig
//der service muss auf alle nachrichten  der leaves waren und kann diese sortieren
type Travers struct {
	Caller     *actor.PID
	TreeValues map[int]string
}
type Search struct {
	Caller *actor.PID
	Key    int
}

type Insert struct {
	Key   int
	Value string
}
type Delete struct {
	Key        int
	Caller     *actor.PID
	NeedUpdate []*actor.PID
}

type Error struct {
	OriginalMsg interface{}
}

type Success struct {
	Key         int
	Value       string
	OriginalMsg interface{}
}
