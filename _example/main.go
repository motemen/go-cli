//go:generate go run _tools/gen.go -out cmds.go $GOFILE

package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/motemen/go-cli"
)

func main() {
	cli.Run(os.Args[1:])
}

/*
+main - greeting

	main -name=motemen

if no subcommands are specified, greet motemen (or specified name by -name flag)
with a bow.
*/
func actionMain(flags *flag.FlagSet, args []string) error {
	var name string
	flags.StringVar(&name, "name", "motemen", "hello!")

	flags.Parse(args)
	log.Printf("hello! %s :bow:\n", name)
	return nil
}

/*
+command up - count up!

	up [-f <from>] <count>

Counts up to specified count. If -f flag was specified, counting starts with
that number.
*/
func actionUp(flags *flag.FlagSet, args []string) error {
	var from int
	flags.IntVar(&from, "f", 1, "count starts from this number")
	flags.Parse(args)

	args = flags.Args()
	if len(args) < 1 {
		return cli.ErrUsage
	}

	count, err := strconv.ParseInt(args[0], 0, 0)
	if err != nil {
		return err
	}

	for i := from; i <= int(count); i++ {
		log.Printf("count: %d\n", i)
	}

	return nil
}

// +command smile - show smile
//
// 	smile
//
// Shows smile.
//
// NOTE: as this action does not call flags.Parse(), passing -h to this command
// does not show the help.
func actionSmile(flags *flag.FlagSet, args []string) error {
	log.Println("( ╹◡╹)")
	return nil
}
