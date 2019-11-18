package service

import (
	"crypto/md5"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/logger"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/tree"
	"time"
)

type Service struct {
	//Place to store all the trees managed by the service
	trees  map[int32]Tree
	nextId func() int32
}

type Tree struct {
	Root  *actor.PID
	Token string
}

//hier werden nur cli MSGs empfangen
func (service *Service) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		logger.GetInstance().Info.Println("started service")
		service.nextId = idGenerator()
		service.trees = make(map[int32]Tree)
	case *messages.Create:
		logger.GetInstance().Info.Println("recieved create")
		service.create(msg, context)
	case *messages.Insert:
		service.insert(msg, context)
	case *messages.Search:
		service.search(msg, context)
	case *messages.Delete:
		service.delete(msg, context)
	case *messages.Traverse:
		service.traverse(msg, context)
	case *messages.Remove:
		service.remove(msg, context)
	}
}

func (service *Service) create(msg *messages.Create, context actor.Context) {
	id := service.nextId()
	token := generateToken()
	root := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return tree.NewRoot(msg.MaxElems)
	}))
	service.trees[id] = Tree{
		Root:  root,
		Token: token,
	}
	msg.Response = &messages.Create_Response{
		Success: true,
		Id:      id,
		Token:   token,
	}
	logger.GetInstance().Info.Printf("Created Tree ID: %v, token: %v root: %v\n", id, token, root)
	context.Respond(msg)
}

func (service *Service) insert(msg *messages.Insert, context actor.Context) {
	root, ok := service.getRootNode(msg.Id, msg.Token)
	if ok {
		helper := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &inserter{
				cli: context.Sender(),
				msg: *msg,
			}
		}))
		context.RequestWithCustomSender(root, &tree.Insert{Key: msg.Key, Value: msg.Value}, helper)
		logger.GetInstance().Info.Printf("Inserter PID : %v\n", helper)
		logger.GetInstance().Info.Printf("Started insert ID:%v,token: %v, root:%v key %v val %v\n", msg.Id, msg.Token, root, msg.Key, msg.Value)
	} else {
		logger.GetInstance().Info.Println("Wrong credentials for insert")
		msg.Response = &messages.Insert_Response{
			Success: false,
			Error:   "Wrong credentials",
		}
		context.Respond(msg)
	}
}

func (service *Service) search(msg *messages.Search, context actor.Context) {
	root, ok := service.getRootNode(msg.Id, msg.Token)
	if ok {
		helper := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &searcher{
				cli: context.Sender(),
				msg: *msg,
			}
		}))
		context.RequestWithCustomSender(root, &tree.Search{Key: msg.Key}, helper)
		logger.GetInstance().Info.Printf("Started search for %v ID:%v,token: %v, root:%v \n", msg.Key, msg.Id, msg.Token, root)
	} else {
		logger.GetInstance().Info.Printf("Wrong credentials for search\n")
		msg.Response = &messages.Search_Response{
			Success: false,
			Error:   "Wrong credentials",
		}
		context.Respond(msg)
	}
}

func (service *Service) delete(msg *messages.Delete, context actor.Context) {
	root, ok := service.getRootNode(msg.Id, msg.Token)
	if ok {
		helper := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &deleter{
				cli: context.Sender(),
				msg: *msg,
			}
		}))
		logger.GetInstance().Info.Printf("Started delete for %v ID:%v,token: %v, root:%v \n", msg.Key, msg.Id, msg.Token, root)
		context.RequestWithCustomSender(root, &tree.Delete{Key: msg.Key}, helper)
	} else {
		logger.GetInstance().Info.Printf("Wrong credentials for delete\n")
		msg.Response = &messages.Delete_Response{
			Success: false,
			Error:   "Wrong credentials",
		}
		context.Respond(msg)
	}
}

func (service *Service) traverse(msg *messages.Traverse, context actor.Context) {
	root, ok := service.getRootNode(msg.Id, msg.Token)
	if ok {
		//spawn traversactor
		traversActorPID := context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &traverser{
				cli:           context.Sender(),
				msg:           *msg,
				nMessagesWait: 1,
				treemap:       make(map[int32]string),
			}
		}))
		logger.GetInstance().Info.Printf("Started traverse ID:%v,token: %v, root:%v \n", msg.Id, msg.Token, root)
		context.RequestWithCustomSender(root, &tree.Traverse{}, traversActorPID)
	} else {
		logger.GetInstance().Info.Printf("Wrong credentials for traverse\n")
		msg.Response = &messages.Traverse_Response{
			Success: false,
			Error:   "Wrong credentials",
		}
		context.Respond(msg)
	}
}
func (service *Service) remove(msg *messages.Remove, context actor.Context) {
	root, ok := service.getRootNode(msg.Id, msg.Token)
	if ok {
		context.Send(root, &tree.Kill{})
		msg.Response = &messages.Remove_Response{
			Success: true,
		}
		context.Respond(msg)
		delete(service.trees, msg.Id)
		logger.GetInstance().Info.Printf("Started remove ID:%v,token: %v, root:%v \n", msg.Id, msg.Token, root)
	} else {
		logger.GetInstance().Info.Printf("Wrong credentials for remove\n")
		msg.Response = &messages.Remove_Response{
			Success: false,
			Error:   "Wrong credentials",
		}
		context.Respond(msg)
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
	h := md5.New()
	token := h.Sum([]byte(t))
	return fmt.Sprintf("%x", token)
}

//method to get the pid of the root of the tree matching the token and id
//if none mach false is returned
func (service *Service) getRootNode(id int32, token string) (*actor.PID, bool) {
	value, ok := service.trees[id]
	if !ok {
		return nil, false
	}
	if value.Token != token {
		return nil, false
	}
	return value.Root, true
}
