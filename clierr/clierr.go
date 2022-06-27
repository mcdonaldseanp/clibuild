package clierr

import "fmt"

type InvalidInput struct {
	Message string
	Origin  error
}

func (ii *InvalidInput) Error() string {
	if ii.Origin != nil {
		return fmt.Sprintf("invalid input\n%s\n\nTrace:\n%s\n", ii.Message, ii.Origin)
	} else {
		return fmt.Sprintf("invalid input\n%s\n", ii.Message)
	}
}
