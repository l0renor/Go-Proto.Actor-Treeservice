package tree

import "github.com/AsynkronIT/protoactor-go/actor"

//Die nachricht verläuft den baum nach unten und teilt sich auf so oft nötig
//der service muss auf alle nachrichten  der leaves waren und kann diese sortieren
type Travers struct {
	Caller     *actor.PID
	TreeValues map[int32]string
}

type Search struct {
	Caller *actor.PID
	Key    int32
}

type Insert struct {
	Key   int32
	Value string
}

type Delete struct {
	Key        int32
	Caller     *actor.PID
	NeedUpdate []*actor.PID
}
type UpdateMaxleft struct {
	NewValue int32
}

type Error struct {
	OriginalMsg interface{}
}

type Success struct {
	Key         int32
	Value       string
	OriginalMsg interface{}
}
