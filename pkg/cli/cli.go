package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func TTYPrompt(prompt, defaultValue string) (string, error) {
	if defaultValue != "" {
		prompt += " [" + defaultValue + "]"
	}

	fmt.Print(prompt + ": ")

	value, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if value == "\n" {
		value = defaultValue
	}

	return strings.TrimSpace(value), err
}

func TTYPassword() (string, error) {
	fmt.Print("password: ")
	defer fmt.Println()

	password, err := term.ReadPassword(int(os.Stdin.Fd()))

	return string(password), err
}

func TTYPin(name string) func() (string, error) {
	return func() (string, error) {
		fmt.Printf("%s PIN: ", name)
		defer fmt.Println()

		pin, err := term.ReadPassword(int(os.Stdin.Fd()))

		return string(pin), err
	}
}
