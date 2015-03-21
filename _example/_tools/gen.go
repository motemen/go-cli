package main

import (
	"log"
	"os"

	"github.com/motemen/cli/gen"
)

func main() {
	out, err := os.Create("cmds.go")
	if err != nil {
		log.Fatal(err)
	}

	err = gen.Generate(out, "main.go", nil)
	if err != nil {
		log.Fatal(err)
	}
}
