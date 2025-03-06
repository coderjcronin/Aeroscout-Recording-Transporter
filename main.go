package main

import (
	"bufio"
	"fmt"
	"os"
)

func prepRecordings() {
	//TODO - Fill this out
}

func prepAnalysis() {
	//TODO - Fill this out
}

func mainMenu() int {
	userInput := bufio.NewReader(os.Stdin)

	fmt.Print("Please select an operation mode:\n\t1) Prepare recordings for transport or storage\n\t2) Prepare recordings for analysis and reporting\n\nPlease type the number followed by the ENTER key (invalid response will exit): ")
	programMode, _ := userInput.ReadString('\n')

	if programMode == "1" {
		return 1
	} else if programMode == "2" {
		return 2
	}

	return 0
}

func main() {

	// Go "while True" loop. Will exit when menuSelection is invalid
	for {
		menuSelection := mainMenu() // Show main menu to select function

		if menuSelection == 1 {
			prepRecordings() // User selected to prep recordings, run function
		} else if menuSelection == 2 {
			prepAnalysis() // User selected to prep analysis, run function
		} else {
			break // Leave the loop
		}
	}
}
