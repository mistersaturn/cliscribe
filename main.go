package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// STYLES
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#3ca4d8")).
		Align(lipgloss.Center).
		Width(45)

	accentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#9C999A")).
		Align(lipgloss.Center).
		Width(45)

	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#D54E53")).
		Align(lipgloss.Center).
		Width(45)
)

// MODEL
type Model struct {
	lines []string
}

// UPDATE
// Handle hotkeys here:
func Update(model *Model, command string) {
	switch {
	case strings.HasPrefix(command, ":C "):
		copyLine(model, command[3:])
	case strings.HasPrefix(command, ":D "):
		deleteLine(model, command[3:])
	case command == ":H":
		displayHelp()
	default:
		if command != ":S" { // Prevent saving ':S' to file
			model.lines = append(model.lines, command)
		}
	}
}

// displayHelp shows the available commands
func displayHelp() {
	fmt.Println("Save: ':S' Copy: ':C <line num>' Delete: ':D <line num>'")
}

// copyLine duplicates a line in the model based on the provided line number
func copyLine(model *Model, lineNumberStr string) {
	lineNumber, err := strconv.Atoi(strings.TrimSpace(lineNumberStr))
	if err != nil || lineNumber < 1 || lineNumber > len(model.lines) {
		fmt.Println(errorStyle.Render("Invalid line number:", lineNumberStr))
		return
	}
	lineToCopy := model.lines[lineNumber-1]
	model.lines = append(model.lines, lineToCopy)
}

// deleteLine removes a line from the model based on the provided line number
func deleteLine(model *Model, lineNumberStr string) {
	lineNumber, err := strconv.Atoi(strings.TrimSpace(lineNumberStr))
	if err != nil || lineNumber < 1 || lineNumber > len(model.lines) {
		fmt.Println(errorStyle.Render("Invalid line number:", lineNumberStr))
		return
	}
	model.lines = append(model.lines[:lineNumber-1], model.lines[lineNumber:]...)
}

// saveToFile writes the model's lines to the specified file
func saveToFile(file *os.File, model *Model) {
	file.Truncate(0)
	file.Seek(0, 0)
	for _, line := range model.lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			fmt.Println(errorStyle.Render("Error saving file."), err)
			return
		}
	}
}

// readLines reads all lines from the file and returns them as a slice
func readLines(file *os.File) ([]string, error) {
	var lines []string
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// VIEW
// printFileContents displays the contents of the model
func printFileContents(model *Model) {
	fmt.Println(accentStyle.Render("----------------------"))
	for _, line := range model.lines {
		fmt.Println(line)
	}
}

// MAIN
func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(titleStyle.Render("cliScribe -- 1.3.0"))
	fmt.Print("\nOpen File\n-> ")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(errorStyle.Render("Error opening file:"), err)
		return
	}
	defer file.Close()

	model := &Model{}
	if model.lines, err = readLines(file); err != nil {
		fmt.Println(errorStyle.Render("Error reading file:"), err)
		return
	}
	printFileContents(model)

	fmt.Println(accentStyle.Render("Enter New Text -> (Type ':H' for Help)"))
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		Update(model, input)
		// Save file:
		if input == ":S" {
			saveToFile(file, model)
			fmt.Println("Saved.")
			break
		}
	}
}
