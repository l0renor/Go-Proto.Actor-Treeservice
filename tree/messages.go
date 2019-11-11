package tree

import "github.com/AsynkronIT/protoactor-go/actor"

//Die Nachricht verläuft den Baum nach unten und teilt sich auf so oft nötig
//der service muss auf alle nachrichten  der leaves waren und kann diese sortieren
type Travers struct {
	Caller     *actor.PID
	TreeValues map[int32]string
}

type Search struct {
	Key int32
}

type Insert struct {
	Key   int32
	Value string
}

type Delete struct {
	Key        int32
	NeedUpdate *actor.PID
}

type UpdateMaxLeft struct {
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
