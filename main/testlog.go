package main

import (
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/logger"
)

func main() {
	log := logger.GetInstance()
	log.Info.Println("info test")
	log.Error.Println("Fail")
}
