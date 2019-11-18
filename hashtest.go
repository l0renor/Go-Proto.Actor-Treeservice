package main

import (
	"crypto/sha1"
	"encoding/hex"
	"time"
)

func main() {
	t := time.Now().String()
	h := sha1.New()
	h.Write([]byte(t))
	token := h.Sum(nil)
	print(hex.Dump(token))
}
