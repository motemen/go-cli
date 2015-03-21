package main

import (
	"flag"
	"log"
	"os"

	"github.com/motemen/go-cli/gen"
)

// go run _tools/gen.go -out generated.go source.go
func main() {
	out := flag.String("out", "", "output file")
	flag.Parse()

	if *out == "" {
		log.Fatal("-out should be specified")
	}

	in := flag.Arg(0)
	if in == "" {
		log.Fatal("input file required")
	}

	w, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}

	err = gen.Generate(w, flag.Arg(0), nil)
	if err != nil {
		log.Fatal(err)
	}
}
