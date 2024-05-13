//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {

	fmt.Println("Press ESC button or Ctrl-C to exit this program")
	fmt.Println("Press any key to see their ASCII code follow by Enter")

	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	for {
		// only read single characters, the rest will be ignored!!
		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		fmt.Print(">")
		input, _ := consoleReader.ReadByte()

		ascii := input

		// ESC = 27 and Ctrl-C = 3
		if ascii == 27 || ascii == 3 {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		fmt.Println("ASCII : ", ascii)
	}

}
