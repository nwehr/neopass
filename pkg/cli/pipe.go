package cli

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func weAreInAPipe() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	if info.Mode()&os.ModeCharDevice == 0 {
		return true
	}

	return false
}

func readFromPipe() (string, error) {
	r := bufio.NewReader(os.Stdin)

	var output []rune

	for {
		input, _, err := r.ReadRune()
		if err != nil && err == io.EOF {
			r.ReadRune()
			break
		}
		output = append(output, input)
	}

	runesToString := func(runes []rune) (outString string) {
		for _, v := range runes {
			outString += string(v)
		}
		return
	}

	name := strings.TrimSpace(runesToString(output))
	return name, nil
}
