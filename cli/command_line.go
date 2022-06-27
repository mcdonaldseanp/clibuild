package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mcdonaldseanp/clibuild/clierr"
	"github.com/mcdonaldseanp/clibuild/clivrsn"
)

type Command struct {
	Verb        string
	Noun        string
	Supports    []string
	ExecutionFn func()
}

// shouldHaveArgs does two things:
// * validate that the number of args that aren't flags have been provided (i.e. the number of strings
//    after the command name that aren't flags)
// * parse the remaining flags
//
// If the wrong number of args is passed it prints helpful usage
func ShouldHaveArgs(num_args int, usage string, description string, flagset *flag.FlagSet) {
	real_args := num_args + 3
	passed_fs := flagset != nil
	for index, arg := range os.Args {
		if arg == "-h" {
			fmt.Fprintf(os.Stderr, "Usage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
			if passed_fs {
				fmt.Fprintf(os.Stderr, "Available flags:\n")
				flagset.PrintDefaults()
			}
			os.Exit(0)
		}
		// None of the arguments required by the command should start with dashes, if they
		// do assume an arg is missing and this is a flag
		if index <= num_args && strings.HasPrefix(arg, "-") {
			fmt.Fprintf(os.Stderr, "Error running command:\n\nInvalid input, not enough arguments.\n\nUsage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
			if passed_fs {
				fmt.Fprintf(os.Stderr, "Available flags:\n")
				flagset.PrintDefaults()
			}
			os.Exit(1)
		}
	}
	if len(os.Args) < real_args {
		fmt.Fprintf(os.Stderr, "Error running command:\n\nInvalid input, not enough arguments.\n\nUsage:\n  %s\n\nDescription:\n  %s\n\n", usage, description)
		if passed_fs {
			fmt.Fprintf(os.Stderr, "Available flags:\n")
			flagset.PrintDefaults()
		}
		os.Exit(1)
	} else if len(os.Args) > real_args && passed_fs {
		flagset.Parse(os.Args[real_args:])
	}
}

func osSupportsCommand(cmd Command) bool {
	for _, os_name := range cmd.Supports {
		if runtime.GOOS == os_name {
			return true
		}
	}
	return false
}

func printTopUsage(tool_name string, command_list []Command) {
	fmt.Printf("Usage:\n  %s [VERB] [NOUN] [ARGUMENTS] [FLAGS]\n\nAvailable commands:\n", tool_name)
	for _, command := range command_list {
		if osSupportsCommand(command) {
			fmt.Printf("    %s %s\n", command.Verb, command.Noun)
		}
	}
}

// handleCommandAirer catches InvalidInput airer.Airers and prints usage
// if that was the error thrown. IF a different type of airer.Airer is thrown
// it just prints the error.
//
// If the command succeeds handleCommandAirer exits the whole go process
// with code 0
func HandleCommandError(err error, usage string, description string, flagset *flag.FlagSet) {
	if err != nil {
		switch err.(type) {
		case *clierr.InvalidInput:
			fmt.Fprintf(os.Stderr, "%s\nUsage:\n  %s\n\nDescription:\n  %s\n\n", err, usage, description)
			if flagset != nil {
				flagset.PrintDefaults()
			}
		default:
			fmt.Fprintf(os.Stderr, "Error running command:\n\n%s\n", err)
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func RunCommand(tool_name string, tool_version string, command_list []Command) {
	default_command_list := []Command{
		{
			Verb:     "update",
			Noun:     "version",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := tool_name + " update version [VERSION FILE] [NEW VERSION]"
				description := "Update " + tool_name + "'s version, defaults to the next Z release if no [NEW VERSION] is given"
				ShouldHaveArgs(1, usage, description, nil)
				new_version := ""
				if len(os.Args) > 4 {
					new_version = os.Args[4]
				}
				HandleCommandError(
					clivrsn.UpdateVersion(os.Args[3], new_version),
					usage,
					description,
					nil,
				)
			},
		},
	}
	command_list = append(command_list, default_command_list...)

	if len(os.Args) > 2 {
		for _, command := range command_list {
			if os.Args[1] == command.Verb && os.Args[2] == command.Noun && osSupportsCommand(command) {
				command.ExecutionFn()
			}
		}
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version":
			fmt.Fprintf(os.Stdout, "%s\n", tool_version)
			os.Exit(0)
		case "-h":
			printTopUsage(tool_name, command_list)
			os.Exit(0)
		}
	}

	// If we've arrived here, that means the args passed don't match an existing command
	// --version or -h
	fmt.Printf("Unknown %s command \"%s\"\n\n", tool_name, strings.Join(os.Args, " "))
	printTopUsage(tool_name, command_list)
	os.Exit(1)
}
