package cli

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

// App represents for a CLI program with commands.
type App struct {
	Name     string
	Commands map[string]*Command
}

// Command represents one of commands of App.
type Command struct {
	Name string

	// Short (one line) description of the command. Used when the program as
	// invoked without a command name
	Short string

	// Long description of the command. The first line of Long should be a
	// usage line i.e. it starts with the command name. Used when invoked with -h
	Long string

	// The actual implementation of the command. The function will receive two
	// arguments, namely flags and args.  flags is an flag.FlagSet and args are
	// the command line arguments after the command name.  flags is not
	// initialized, so declaring flag variables and arguments parsing with
	// flags.Parse(args) should be called inside the function.
	//
	// Return ErrUsage if you want to show user the command usage.
	Action func(flags *flag.FlagSet, args []string) error
}

var (
	// Default implementation of App. Its name is set to os.Args[0].
	Default = &App{
		Commands: Commands,
	}
	// Commands is default value of Default.Commands.
	Commands = map[string]*Command{}
)

// Run is a shortcut for Default.Run.
func Run(args []string) { Default.Run(args) }

// Use is a shortcut for Default.Use.
func Use(cmd *Command) { Default.Use(cmd) }

// ErrUsage is the error indicating the user had wrong usage.
var ErrUsage = fmt.Errorf("usage error")

func init() {
	Default.Name = os.Args[0]
}

// Run is the entry point of the program. It recognizes the first element of
// args as a command name, and dispatches a command with rest arguments.
func (app *App) Run(args []string) {
	if len(args) == 0 {
		app.PrintUsage()
		os.Exit(2)
	}

	cmdName := args[0]
	if cmd, ok := app.Commands[cmdName]; ok {
		flags := flag.NewFlagSet(cmdName, flag.ExitOnError)
		flags.Usage = func() {
			fmt.Fprintln(os.Stderr, cmd.Usage(flags))
		}

		err := cmd.Action(flags, args[1:])
		if err != nil {
			if err == ErrUsage {
				flags.Usage()
				os.Exit(2)
			} else {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	} else {
		app.PrintUsage()
		os.Exit(2)
	}
}

// PrintUsage prints out the usage of the program with its commands listed.
func (app *App) PrintUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <command> [<args>]\n\n", app.Name)
	fmt.Fprintf(os.Stderr, "Commands:\n")

	names := make([]string, 0, len(app.Commands))
	for name := range app.Commands {
		names = append(names, name)
	}

	sort.Strings(names)

	w := tabwriter.NewWriter(os.Stderr, 0, 8, 4, ' ', 0)
	for _, name := range names {
		fmt.Fprintf(w, "    %s\t%s\n", name, app.Commands[name].Short)
	}
	w.Flush()
}

// Use registers a app command cmd.
func (app *App) Use(cmd *Command) {
	app.Commands[cmd.Name] = cmd
}

// Usage returns a usage documentation of a command.
func (c Command) Usage(flags *flag.FlagSet) string {
	usage := fmt.Sprintf("Usage: %s", c.Long)

	if flags == nil || flags.NFlag() == 0 {
		return usage
	}

	buf := bytes.NewBufferString(usage)
	buf.WriteString("\n\nOptions:\n")

	defer flags.SetOutput(nil)
	flags.SetOutput(buf)

	flags.PrintDefaults()

	return buf.String()
}
