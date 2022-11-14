package main

import (
	"math/rand"
	"qq/client/cmd"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	cmd.Execute()
}
