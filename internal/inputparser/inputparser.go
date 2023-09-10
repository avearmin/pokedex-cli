package inputparser

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type inputParser struct {
	scanner *bufio.Scanner
	args    []string
	len     int
}

func NewInputParser(len int) *inputParser {
	scanner := bufio.NewScanner(os.Stdin)
	return &inputParser{
		scanner: scanner,
		args:    make([]string, len),
		len:     len,
	}
}

func (in *inputParser) Arg(index int) string {
	return in.args[index]
}

func (in *inputParser) ScanAndParse() error {
	in.scanner.Scan()
	input := in.scanner.Text()
	input = trimAndLower(input)
	if err := in.parse(input); err != nil {
		return err
	}
	return nil
}

func (in *inputParser) parse(input string) error {
	args := strings.Fields(input)
	if in.len < len(args) {
		return errors.New("Input has too many args")
	}
	for i := 0; i < in.len; i++ {
		if i < len(args) {
			in.args[i] = args[i]
		} else {
			in.args[i] = ""
		}
	}
	return nil
}

func trimAndLower(s string) string {
	newString := strings.ToLower(s)
	newString = strings.Trim(newString, " ")
	return newString
}
