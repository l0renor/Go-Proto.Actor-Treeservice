package cli

import (
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	l "github.com/AsynkronIT/protoactor-go/log"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/messages"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// nolint:gochecknoglobals
var bindAddr string

// nolint:gochecknoglobals
var remoteAddr string

// nolint:gochecknoglobals
var id int

// nolint:gochecknoglobals
var token string

// nolint:gocognit
func Main() {
	app := &cli.App{
		Name:  "treecli",
		Usage: "Interact with a specified treeservice",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "bind",
				Usage:       "Address `HOST:PORT` to bind CLI to",
				Value:       "localhost:8091",
				Destination: &bindAddr,
			},
			&cli.StringFlag{
				Name:        "remote",
				Usage:       "Address `HOST:PORT` of remote service",
				Value:       "localhost:8090",
				Destination: &remoteAddr,
			},
			&cli.IntFlag{
				Name:        "id",
				Usage:       "ID of target Tree",
				Destination: &id,
			},
			&cli.StringFlag{
				Name:        "token",
				Usage:       "Token of target Tree",
				Destination: &token,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Creates a new tree",
				Action: func(context *cli.Context) error {
					if context.NArg() == 1 && id == 0 && token == "" {
						maxElems, err := strconv.Atoi(context.Args().Get(0))
						if err != nil {
							return errors.New("first argument has to be an integer")
						}
						msg := &messages.Create{MaxElems: int32(maxElems)}
						err = callService(msg)
						return err
					}
					return errors.New("call create with one argument and no credential flags")
				},
			},
			{
				Name:  "insert",
				Usage: "Inserts key-value pair",
				Action: func(context *cli.Context) error {
					if context.NArg() == 2 && id != 0 && token != "" {
						key, err := strconv.Atoi(context.Args().Get(0))
						if err != nil {
							return errors.New("first argument has to be an integer")
						}
						value := context.Args().Get(1)
						msg := &messages.Insert{
							Id:    int32(id),
							Token: token,
							Key:   int32(key),
							Value: value,
						}
						err = callService(msg)
						return err
					}
					return errors.New("call insert with two arguments and credential flags")
				},
			},
			{
				Name:  "search",
				Usage: "Searches by key",
				Action: func(context *cli.Context) error {
					if context.NArg() == 1 && id != 0 && token != "" {
						key, err := strconv.Atoi(context.Args().Get(0))
						if err != nil {
							return errors.New("first argument has to be an integer")
						}
						msg := &messages.Search{
							Id:    int32(id),
							Token: token,
							Key:   int32(key),
						}
						err = callService(msg)
						return err
					}
					return errors.New("call search with one argument and credential flags")
				},
			},
			{
				Name:  "delete",
				Usage: "Deletes one key-value pair",
				Action: func(context *cli.Context) error {
					if context.NArg() == 1 && id != 0 && token != "" {
						key, err := strconv.Atoi(context.Args().Get(0))
						if err != nil {
							return errors.New("first argument has to be an integer")
						}
						msg := &messages.Delete{
							Id:    int32(id),
							Token: token,
							Key:   int32(key),
						}
						err = callService(msg)
						return err
					}
					return errors.New("call delete with one argument and credential flags")
				},
			},
			{
				Name:  "traverse",
				Usage: "Traverses trough whole tree",
				Action: func(context *cli.Context) error {
					if context.NArg() == 0 && id != 0 && token != "" {
						msg := &messages.Traverse{
							Id:    int32(id),
							Token: token,
						}
						err := callService(msg)
						return err
					}
					return errors.New("call traverse with credential flags and no arguments")
				},
			},
			{
				Name:  "remove",
				Usage: "Removes full tree",
				Action: func(context *cli.Context) error {
					if context.NArg() == 0 && id != 0 && token != "" {
						msg := &messages.Remove{
							Id:    int32(id),
							Token: token,
						}
						err := callService(msg)
						return err
					}
					return errors.New("call remove with credential flags and no arguments")
				},
			},
		},
	}
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		os.Exit(0)
	}()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func callService(msg interface{}) error {
	remote.Start(bindAddr)
	var wg sync.WaitGroup
	rootContext := actor.EmptyRootContext
	actor.SetLogLevel(l.ErrorLevel)
	receiver := rootContext.Spawn(actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &receiver{wg: &wg}
	}))
	response, err := remote.SpawnNamed(remoteAddr, "remote", "tree", 5*time.Second)
	if err != nil {
		return err
	}
	service := response.Pid
	rootContext.RequestWithCustomSender(service, msg, receiver)
	wg.Wait()
	return nil
}
