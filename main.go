package main

import (
	"fmt"
	"os"

	"github.com/JanGordon/melte/cli"
)

func main() {
	// fmt.Println("HEllo World")
	// compile.Transform()
	c := 0
	os.Create(fmt.Sprintf("da%d.js", c))
	startCli()
}

func startCli() {
	cli.Start()
}
