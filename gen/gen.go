// Package gen generates github.com/motemen/go-cli.Command from
// function docs in files specified.
package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"regexp"
	"strings"
)

const fileFormat = `// auto-generated file

package main

import "github.com/motemen/go-cli"

func init() {%s}
`

const commandFormat = `
	cli.Use(
		&cli.Command{
			Name:   %q,
			Action: %s,
			Short:  %q,
			Long:   %q,
		},
	)
`

var commandReg = regexp.MustCompile(`^ *(\S+) +- +(.+)\n((?s).+)`)
var mainReg = regexp.MustCompile(`^ *- +(.+)\n((?s).+)`)

/*
Generate reads source file for command actions with their usage documentations
and writes Go code that registers the command to cli.

Usage documentation should be like below:

	// +command <name> - <short>
	//
	// <usage line>
	//
	// <long description>...
	func action(flags *flag.FlagSet, args []string) {
	}

Currently, generated files will be like below:

	// auto-generated file

	package main

	import "github.com/motemen/go-cli"

	func init() {
	    cli.Use(
	        &cli.Command{
	            Name:   "foo",
	            ...
	        }
	    )
	    ...
	}

You can define `main action` without sub-command.
Usage documentation of main action should be like below:

	// +main - <short>
	//
	// <usage line>
	//
	// <long description>...
	func mainAction(flags *flag.FlagSet, args []string) {
		...
	}
*/
func Generate(w io.Writer, path string, src interface{}) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return err
	}

	commandCodes := []string{}

	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Doc == nil {
			continue
		}

		doc := funcDecl.Doc.Text()
		isMain := false
		pos := strings.Index(doc, "+command")
		if pos == -1 {
			pos = strings.Index(doc, "+main")
			if pos == -1 {
				continue
			}
			isMain = true
		}

		var (
			name  string
			short string
			long  string
		)
		if isMain {
			doc = doc[pos+len("+main "):]
			m := mainReg.FindStringSubmatch(doc)
			if m == nil {
				continue
			}
			short = m[1]
			long = strings.TrimSpace(m[2])
		} else {
			doc = doc[pos+len("+command "):]
			m := commandReg.FindStringSubmatch(doc)
			if m == nil {
				continue
			}
			name = m[1]
			short = m[2]
			long = strings.TrimSpace(m[3])
		}

		commandCodes = append(
			commandCodes,
			fmt.Sprintf(commandFormat, name, funcDecl.Name.Name, short, long),
		)
	}

	code := fmt.Sprintf(fileFormat, strings.Join(commandCodes, ""))

	_, err = w.Write([]byte(code))
	return err
}
