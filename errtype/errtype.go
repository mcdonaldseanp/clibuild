package errtype

import "fmt"

// Invalid input errors represent user input errors
//
// mostly returned when validators fail validation of inputs from the CLI
type InvalidInput struct {
	Message string
	Origin  error
}

func (ii *InvalidInput) Error() string {
	if ii.Origin != nil {
		return fmt.Sprintf("invalid input\n%s\n\ntrace:\n%s\n", ii.Message, ii.Origin)
	} else {
		return fmt.Sprintf("invalid input\n%s\n", ii.Message)
	}
}

// Shell errors represent failures to run shell commands on localhost
type ShellError struct {
	Message string
	Origin  error
}

func (se *ShellError) Error() string {
	if se.Origin != nil {
		return fmt.Sprintf("shell execution failed\n%s\n\ntrace:\n%s\n", se.Message, se.Origin)
	} else {
		return fmt.Sprintf("shell execution failed\n%s\n", se.Message)
	}
}

// Remote Shell errors represent failures running commands on remote targets
type RemoteShellError struct {
	Message string
	Origin  error
}

func (rs *RemoteShellError) Error() string {
	if rs.Origin != nil {
		return fmt.Sprintf("shell execution failed\n%s\n\ntrace:\n%s\n", rs.Message, rs.Origin)
	} else {
		return fmt.Sprintf("shell execution failed\n%s\n", rs.Message)
	}
}
