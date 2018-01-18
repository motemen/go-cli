package cli

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var exitCode int

func init() {
	exit = func(code int) {
		exitCode = code
	}
}

func TestDefault_Use_Run(t *testing.T) {
	var fooArgs []string
	cmdFoo := &Command{
		Name: "foo",
		Action: func(flags *flag.FlagSet, args []string) error {
			fooArgs = args
			return nil
		},
	}

	Use(cmdFoo)

	if Commands["foo"] != cmdFoo {
		t.Fatal("command should be registered")
	}

	exitCode = 0
	Run([]string{"foo", "x", "y"})

	if len(fooArgs) != 2 || fooArgs[0] != "x" || fooArgs[1] != "y" {
		t.Fatalf("should receive arguments: %v", fooArgs)
	}

	if exitCode != 0 {
		t.Fatalf("exit code must be 0, got %v", exitCode)
	}
}

func TestUsage(t *testing.T) {
	var out bytes.Buffer

	cmd1 := &Command{
		Name:  "cmd1",
		Short: "blah blah blah",
		Long: `cmd1

blah
blah
blah
`,
		Action: func(flags *flag.FlagSet, args []string) error {
			return ErrUsage
		},
	}

	cmd2 := &Command{
		Name:  "cmd2",
		Short: "xyz",
		Long:  "cmd2 [-v] -from <from> -to <to>",
		Action: func(flags *flag.FlagSet, args []string) error {
			flags.Bool("v", false, "set verbosity")
			flags.String("from", "", "specify from")
			flags.String("to", "", "specify to")
			return flags.Parse(args)
		},
	}

	cmd3 := &Command{
		Name: "cmd3",
		Action: func(flags *flag.FlagSet, args []string) error {
			return fmt.Errorf("internal error")
		},
	}

	app := &App{
		Name:        "prog",
		Commands:    make(map[string]*Command),
		ErrorWriter: &out,
	}
	app.Use(cmd1)
	app.Use(cmd2)
	app.Use(cmd3)

	var s string

	runOut := func(args ...string) string {
		out.Reset()
		app.Run(args)
		return out.String()
	}

	// prog
	s = runOut()
	if strings.HasPrefix(s, "Usage: prog <command> [<args>]\n") == false {
		t.Errorf("should begin with program usage line:\n%v", s)
	}
	if strings.Contains(s, "cmd1") == false {
		t.Error("should contain command name cmd1")
	}
	if strings.Contains(s, "cmd2") == false {
		t.Error("should contain command name cmd2")
	}

	// prog cmd1
	s = runOut("cmd1")
	if strings.Contains(s, cmd1.Long) == false {
		t.Errorf("should include long description:\n%v", s)
	}

	// prog cmd2 -h
	s = runOut("cmd2", "-h")
	if strings.Contains(s, "set verbosity") == false {
		t.Errorf("should include flag specs:\n%v", s)
	}

	// prog cmd3
	s = runOut("cmd3")
	if strings.Contains(s, "internal error") == false {
		t.Errorf("should include error inside action:\n%v", s)
	}

	// prog cmdX
	s = runOut("cmdX", "-h")
	if strings.HasPrefix(s, "Usage: prog <command> [<args>]\n") == false {
		t.Errorf("should begin with program usage line:\n%v", s)
	}
}

func TestCommand_Usage(t *testing.T) {
	cmd := &Command{
		Name:  "cmd",
		Short: "short usage",
		Long:  "cmd long usage...",
		Action: func(flags *flag.FlagSet, args []string) error {
			return nil
		},
	}

	usageWithoutFlags := cmd.Usage(nil)
	if usageWithoutFlags != "Usage: cmd long usage..." {
		t.Errorf("unexpected usage:\n%v", usageWithoutFlags)
	}
}

func TestApp_resolveCmd(t *testing.T) {
	cmd0 := &Command{
		Name: "",
	}
	cmd1 := &Command{
		Name: "cmd1",
	}
	cmd2 := &Command{
		Name: "cmd2",
	}
	cmd3 := &Command{
		Name: "cmd3",
	}
	app := &App{
		Name:     "prog",
		Commands: make(map[string]*Command),
	}
	app.Use(cmd0)
	app.Use(cmd1)
	app.Use(cmd2)
	app.Use(cmd3)

	testCases := []struct {
		name string
		argv []string
		cmd  string
		args []string
	}{
		{
			name: "cmd1 with args",
			argv: []string{"cmd1", "11", "22"},
			cmd:  "cmd1",
			args: []string{"11", "22"},
		},
		{
			name: "empty argv",
			argv: []string{},
			cmd:  "",
			args: []string{},
		},
		{
			name: "main with args",
			argv: []string{"arg1", "arg2"},
			cmd:  "",
			args: []string{"arg1", "arg2"},
		},
		{
			name: "main with help",
			argv: []string{"-h"},
			cmd:  "",
			args: []string{"-h"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, args := app.resolveCmd(tc.argv)
			if cmd.Name != tc.cmd {
				t.Errorf("command name miss matched. out: %q, expected name: %q", cmd.Name, tc.cmd)
			}
			if !reflect.DeepEqual(tc.args, args) {
				t.Errorf("args miss matched.\n   out: %#v\nexpect: %#v", args, tc.args)
			}
		})
	}
}
