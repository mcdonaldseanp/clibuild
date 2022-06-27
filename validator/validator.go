package validator

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
)

// Validators should take a string that looks something like:
//
// {"name":"input_param","value":"some_value_to_validate","validate":["NotEmpty","IsNumber"]}

type Validator struct {
	Name     string   `json:"name"`
	Value    string   `json:"value"`
	Validate []string `json:"validate"`
}

func ValidateParams(params string) error {
	unmarshald_data := []Validator{}
	// YAML would have been more human readable, but the use of whitespace in
	// yaml as an actual separator makes it really ugly to create docstrings
	// with yaml in them (which is the primary way this is meant to be used:
	// you pass this function a json docstring that identifies what to
	// validate
	err := json.Unmarshal([]byte(params), &unmarshald_data)
	if err != nil {
		return fmt.Errorf("failed to parse validator as yaml:\n%s", err)
	}
	for _, data := range unmarshald_data {
		for _, validate_type := range data.Validate {
			switch validate_type {
			case "NotEmpty":
				if !(len(data.Value) > 0) {
					return fmt.Errorf("'%s' is empty", data.Name)
				}
			case "IsNumber":
				matcher, _ := regexp.Compile(`^[\d]+$`)
				if !matcher.Match([]byte(data.Value)) {
					return fmt.Errorf("'%s' is not a number, given %s", data.Name, data.Value)
				}
			case "IsIP":
				matcher, _ := regexp.Compile(`^[\d\.]+$`)
				if !matcher.Match([]byte(data.Value)) {
					return fmt.Errorf("'%s' is not an IP address, given %s", data.Name, data.Value)
				}
			case "IsFile":
				files, err := filepath.Glob(data.Value)
				if err != nil {
					return fmt.Errorf("failed attempting to check if '%s' is a file or directory, failure:\n%s", data.Name, err)
				}
				if len(files) < 1 {
					return fmt.Errorf("'%s' is not a file or directory, given %s", data.Name, data.Value)
				}
			default:
				return fmt.Errorf("unknown matcher: %s", validate_type)
			}
		}
	}
	return nil
}
