package clivrsn

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mcdonaldseanp/clibuild/validator"
)

func readVersion(version_file string) (string, error) {
	raw_bytes, arr := readFileInChunks(version_file)
	if arr != nil {
		return "", arr
	}
	lines := strings.Split(string(raw_bytes), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "const VERSION string") {
			split_line := strings.Split(line, " ")
			ver := split_line[len(split_line)-1]
			ver = ver[1 : len(ver)-1]
			return ver, nil
		}
	}
	return "", errors.New("could not find version")
}

func overwriteFile(location string, data []byte) error {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file:\n%s", err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write file:\n%s", err)
	}
	return nil
}

func readFileInChunks(location string) ([]byte, error) {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to open file:\n%s", err)
	}
	defer f.Close()

	// Create a buffer, read 32 bytes at a time
	byte_buffer := make([]byte, 32)
	file_contents := make([]byte, 0)
	for {
		bytes_read, err := f.Read(byte_buffer)
		if bytes_read > 0 {
			file_contents = append(file_contents, byte_buffer[:bytes_read]...)
		}
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("Failed to read file:\n%s", err)
			} else {
				break
			}
		}
	}
	return file_contents, nil
}

func UpdateVersion(version_file string, new_version string) error {
	err := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"version_file","value":"%s","validate":["NotEmpty", "IsFile"]},
			{"name":"new_version","value":"%s","validate":["NotEmpty"]}
		 ]`,
		version_file,
		new_version,
	))
	if err != nil {
		return err
	}
	raw_bytes, arr := readFileInChunks(version_file)
	if arr != nil {
		return arr
	}
	lines := strings.Split(string(raw_bytes), "\n")
	var result string
	for index, line := range lines {
		if strings.HasPrefix(line, "const VERSION string") {
			result = result + fmt.Sprintf("const VERSION string = \"%s\"\n", new_version)
		} else {
			// Don't allow any newlines toward the end of the file
			//
			// This avoids creating a new newline every time the command is run
			// and making more and more newlines the more the command is run
			if len(line) > 0 || index < len(lines)-2 {
				result = result + line + "\n"
			}
		}
	}
	arr = overwriteFile(version_file, []byte(result))
	if arr != nil {
		return arr
	}
	return nil
}

func ReadNextZ(version_file string) (string, error) {
	err := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"version_file","value":"%s","validate":["NotEmpty", "IsFile"]}
		 ]`,
		version_file,
	))
	if err != nil {
		return "", err
	}
	old_version, arr := readVersion(version_file)
	if arr != nil {
		return "", arr
	}
	split_ver := strings.Split(old_version, ".")
	z_release, err := strconv.Atoi(split_ver[2])
	if err != nil {
		return "", fmt.Errorf("could not read next Z version, atoi conversion failed: %s", err)
	}
	z_release++
	split_ver[2] = strconv.Itoa(z_release)
	return strings.Join(split_ver, "."), nil
}
