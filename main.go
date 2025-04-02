package main

// IMPORTS

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//MAIN

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n๛ cliScribe -- 0 . 0 . 1\n")
	fmt.Print("\nOpen File\n->")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fmt.Println("----------------------")
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	fmt.Println("\nEnter New Text -> (Type ':S!' to Save and Exit)")
	var input string
	for {
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == ":S!" {
			break
		}
		_, err := file.WriteString(input + "\n")
		if err != nil {
			fmt.Println(err, err)
		}
	}

	fmt.Println("Saved.")
}
