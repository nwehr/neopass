package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
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

	tty, err := os.Open("/dev/tty")
	if err != nil {
		return "", err
	}

	defer tty.Close()
	defer fmt.Println()

	password, err := terminal.ReadPassword(int(tty.Fd()))

	return string(password), err
}

func TTYPin(name string) func() (string, error) {
	return func() (string, error) {
		fmt.Printf("%s PIN: ", name)

		tty, err := os.Open("/dev/tty")
		if err != nil {
			return "", err
		}

		defer tty.Close()
		defer fmt.Println()

		pin, err := terminal.ReadPassword(int(tty.Fd()))

		return string(pin), err
	}
}
