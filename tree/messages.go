package tree

//Die Nachricht verläuft den Baum nach unten und teilt sich auf so oft nötig
//der service muss auf alle nachrichten  der leaves waren und kann diese sortieren
type Travers struct {
	TreeValues map[int32]string
}
type TraversWaitOneMore struct {
}

type TraverseActor_Msg struct {
	tree []tuple
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
