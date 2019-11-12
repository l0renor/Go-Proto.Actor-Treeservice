package service

import (
	"crypto/sha1"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/tree"
	"time"
)

type Service struct {
	//Place to store all the trees managed by the service
	trees map[int32]Tree
	idGen func() int32
}

type Tree struct {
	Root  *actor.PID
	Token string
}

//hier werden nur CLI MSGs empfangen
func (service *Service) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.Traverse:
		//spawn traversactor
		traversActorPID := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &Traversaktor{
				CLI:           context.Sender(),
				NMessagesWait: 1,
			}
		}))
		//start travers by sending to first node with custom sender = traversactor
		context.RequestWithCustomSender(PID, tree.Travers{}, traversActorPID) //TODO PID des root des trees mit der gewünschten id
	case messages.Insert:
		InsertActorPID := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &InsertActor{CLI: context.Sender()}
		}))
		context.RequestWithCustomSender(PID, tree.Insert{}, InsertActorPID) //TODO PID des root des trees mit der gewünschten id

	case messages.Create:
		id := idGenerator()
		token := generateToken()
		root := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &tree.Node{
				maxElems: node.maxElems,
				inner:    nil,
				leaf:     &Leaf{values: make(map[int32]string, node.maxElems)},
			}
		}))
		service.trees[id] = Tree{
			Root:  nil,
			Token: "",
		}

	case messages.Delete:
		//send tree.kill msg to root
		context.Send()

	}
}

func idGenerator() func() int32 {
	i := 0
	return func() int32 {
		i++
		return int32(i)
	}
}

func generateToken() string {
	t := time.Now().String()
	h := sha1.New()
	h.Write([]byte(t))
	token := h.Sum(nil)
	return string(token)
}
