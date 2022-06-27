package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/mcdonaldseanp/clibuild/cli"
	"github.com/mcdonaldseanp/clibuild/validator"
	"github.com/mcdonaldseanp/clibuild/version"
)

func exampleCommand(input string, example_flag string) error {
	// ValidateParams returns a specific type of error that
	// will produce user friendly output when used in conjunction
	// with HandleCommandError (see below in the main() function)
	err := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"input","value":"%s","validate":["NotEmpty"]}
		 ]`,
		input,
	))
	if err != nil {
		return err
	}
	if input == "show-error" {
		return errors.New("this is what an execution error looks like from the CLI")
	}
	fmt.Printf("Input: %s\nFlags: %s\n", input, example_flag)
	return nil
}

func main() {
	example_fs := flag.NewFlagSet("example_flag_set", flag.ExitOnError)
	example_flag := example_fs.String("example", os.Getenv("EXAMPLE"), "Example flag cli option")

	// All CLI commands should follow naming rules of powershell approved verbs:
	// https://docs.microsoft.com/en-us/powershell/scripting/developer/cmdlet/approved-verbs-for-windows-powershell-commands?view=powershell-7.2
	//
	// The command_list defines the commands available when running from the CLI. Some default commands are added to this list
	// and are available in every project that uses clibuild
	command_list := []cli.Command{
		{
			// Each command should be a combination of verb-then-noun
			Verb: "run",
			Noun: "example",
			// Sometimes commands only apply when running one type (or a subset of types) of operating system. You must
			// specify which OSes are supported for each command
			Supports: []string{"linux", "windows"},
			// The ExecutionFn is what actually runs when a user calls this command from the CLI.
			ExecutionFn: func() {
				usage := "clibuild run example [INPUT] [FLAGS]"
				description := "Example CLI command. Demonstrates an error when [INPUT] is 'show-error'"
				// ShouldHaveArgs is a convenience method to check if the correct number of non-flag args have been passed.
				// "non-flag args" are defined as arguments _after_ the verb and noun that are not flags. In the case of
				// this command there is only 1 non-flag arg "INPUT", everything else is either the verb, noun, or a flag
				//
				// ShouldHaveArgs will print a nicely formatted error response if not enough args were passed on the CLI.
				cli.ShouldHaveArgs(1, usage, description, example_fs)
				// HandleCommandError is another convenience method that checks the error type and returns a well formatted
				// string on stderr if there was an error. This is particularly useful if using a clibuild.validator on the args
				// passed from the CLI, since validators return a specific error that prints a specific user friendly message
				cli.HandleCommandError(
					exampleCommand(os.Args[3], *example_flag),
					usage,
					description,
					example_fs,
				)
			},
		},
	}

	cli.RunCommand("clibuild", version.VERSION, command_list)
}
