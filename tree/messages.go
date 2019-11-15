package tree

//Die Nachricht verläuft den Baum nach unten und teilt sich auf so oft nötig
//der CLI muss auf alle nachrichten  der leaves waren und kann diese sortieren

type Kill struct {
}

type Traverse struct {
	TreeValues map[int32]string
}
type TraverseWaitOneMore struct {
}

type TraverseActor_Msg struct {
	tree []Tuple
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
}

type Success struct {
	Key   int32
	Value string
}

type Tuple struct {
	Key   int32
	Value string
}
