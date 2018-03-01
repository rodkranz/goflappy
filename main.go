package main

import (
	"fmt"
	"os"

	"flappy/game"
)

func main() {
	if err := game.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
	os.Exit(0)
}
