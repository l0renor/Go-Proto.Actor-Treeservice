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
	trees     map[int32]Tree
	nextId    func() int32
	nextToken func() string
}

type Tree struct {
	Root  *actor.PID
	Token string
}

//hier werden nur CLI MSGs empfangen
func (service *Service) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Traverse:
		service.traverse(msg, context)
	case *messages.Insert:
		service.insert(msg, context)

	case *messages.Create:
		id := service.nextId()
		token := service.nextToken()
		root := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &tree.Node{
				MaxElems: msg.MaxElems,
				Inner:    nil,
				Leaf:     &tree.Leaf{values: make(map[int32]string, msg.MaxElems)},
			}
		}))
		service.trees[id] = Tree{
			Root:  root,
			Token: token,
		}

	case *messages.Remove:
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

func (service *Service) insert(msg *messages.Insert, context actor.Context) {
	id := msg.Id
	token := msg.Token
	root, ok := service.getRootNode(id, token)
	if ok {
		InsertActorPID := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &InsertActor{CLI: context.Sender()}
		}))
		context.RequestWithCustomSender(root, tree.Insert{}, InsertActorPID)
	} else {
		//TODO Error wrong token/id
	}
}

func (service *Service) traverse(msg *messages.Traverse, context actor.Context) {
	id := msg.Id
	token := msg.Token
	root, ok := service.getRootNode(id, token)
	if ok {
		//spawn traversactor
		traversActorPID := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &Traversaktor{
				CLI:           context.Sender(),
				NMessagesWait: 1,
			}
		}))
		context.RequestWithCustomSender(root, tree.Travers{}, traversActorPID)
	} else {
		//TODO Error wrong token/id
	}
}

func (service *Service) remove(msg *messages.Remove, context actor.Context) {
	id := msg.Id
	token := msg.Token
	root, ok := service.getRootNode(id, token)
	if ok {
		context.Send(root, tree.Kill{})
	} else {
		//TODO Error wrong token/id
	}
}

//method to get the pid of the root of the tree matching the token and id
//if none mach false is returned
func (service *Service) getRootNode(id int32, token string) (*actor.PID, bool) {
	tree, ok := service.trees[id]
	if !ok {
		return nil, false
	}
	if tree.Token != token {
		return nil, false
	}
	return tree.Root, true
}
