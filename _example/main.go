//go:generate go run _tools/gen.go

package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/motemen/cli"
)

func main() {
	cli.Run(os.Args[1:])
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
