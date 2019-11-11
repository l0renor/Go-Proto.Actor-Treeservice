package tree

import "github.com/ob-vss-ws19/blatt-3-chupa-chups/treeservice"

//Die Nachricht verläuft den Baum nach unten und teilt sich auf so oft nötig
//der CLI muss auf alle nachrichten  der leaves waren und kann diese sortieren
type Travers struct {
	TreeValues map[int32]string
}
type TraversWaitOneMore struct {
}
type Message interface{}
type TraverseActor_Msg struct {
	tree []service.Tuple
}

type Search struct {
	Key int32
}

type Insert struct {
	Key   int32
	Value string
}

type Delete struct {
	Key int32
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
