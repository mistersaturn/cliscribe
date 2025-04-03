package main

// IMPORTS

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MAIN

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n๛ cliScribe -- 1.2.0\n")
	fmt.Print("\nOpen File\n-> ")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	printFileContents(file)

	fmt.Println("\nEnter New Text -> (Type ':H' for Help)\n")
	handleUserInput(file, reader)

	fmt.Println("Saved.")
}

// printFileContents reads and prints the contents of the file
func printFileContents(file *os.File) {
	file.Seek(0, 0) // Reset file pointer to the beginning
	scanner := bufio.NewScanner(file)
	fmt.Println("----------------------")
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// handleUserInput processes user input for adding text and managing lines
func handleUserInput(file *os.File, reader *bufio.Reader) {
	var input string
	for {
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == ":S" {
			break
		} else if strings.HasPrefix(input, ":C ") {
			copyLine(file, input[3:])
		} else if strings.HasPrefix(input, ":D ") {
			deleteLine(file, input[3:])
		} else if input == ":H" {
			fmt.Println("Save: ':S' Copy: ':C <line num>' Delete: ':D <line num>'")
		} else {
			_, err := file.WriteString(input + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}
		}
	}
}

// copyLine duplicates a line in the file based on the provided line number
func copyLine(file *os.File, lineNumberStr string) {
	lineNumber, err := strconv.Atoi(strings.TrimSpace(lineNumberStr))
	if err != nil {
		fmt.Println("Invalid line number:", lineNumberStr)
		return
	}

	lines, err := readLines(file)
	if err != nil {
		fmt.Println("Error reading lines:", err)
		return
	}

	if lineNumber < 1 || lineNumber > len(lines) {
		fmt.Println("Line number out of range.")
		return
	}

	lineToCopy := lines[lineNumber-1]
	_, err = file.WriteString(lineToCopy + "\n")
	if err != nil {
		fmt.Println("Error writing copied line:", err)
	}
}

// deleteLine removes a line from the file based on the provided line number
func deleteLine(file *os.File, lineNumberStr string) {
	lineNumber, err := strconv.Atoi(strings.TrimSpace(lineNumberStr))
	if err != nil {
		fmt.Println("Invalid line number:", lineNumberStr)
		return
	}

	lines, err := readLines(file)
	if err != nil {
		fmt.Println("Error reading lines:", err)
		return
	}

	if lineNumber < 1 || lineNumber > len(lines) {
		fmt.Println("Line number out of range.")
		return
	}

	lines = append(lines[:lineNumber-1], lines[lineNumber:]...) // Remove the line
	file.Truncate(0) // Clear the file
	file.Seek(0, 0)  // Reset file pointer to the beginning
	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing remaining lines:", err)
		}
	}
}

// readLines reads all lines from the file and returns them as a slice
func readLines(file *os.File) ([]string, error) {
	var lines []string
	file.Seek(0, 0) // Reset file pointer to the beginning
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
