package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func prepRecordings() {
	//TODO - Fill this out
}

func prepAnalysis() {
	//TODO - Fill this out
}

/*
	mainMenu - Returns int based on selection

Simply displays the utility main menu and returns an int indicating the user's selection.
1 - Process recordings for transport
2 - Process recordings for analysis
Anything else - Exit
*/
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

/*
	sanityCheck() - returns boolean based on environment values

sanityCheck will check that:
1. We're in a recording session folder (there's a child folder with MAPPlaybackFile.dat)
2. That we can create a file
3. That we can delete the file we created

If we fail at any point, log.Fatal() is raised which will cause an os.Exit() and terminate the program
*/
func sanityCheck() {

	isThereRecording := false

	entries, err := os.ReadDir("./") // Pull the file listing for the local directory

	// If there was an error, log (well, print) it and return false
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the listing for a directory
	for _, file := range entries {
		if file.IsDir() { // Found a directory, let's do a write and delete check, then look for our file!
			tempPath := "./" + file.Name()
			tempFile := tempPath + "/ART.tst"

			f, err := os.Create(tempFile) // Create the temp file, if there's an error we failed
			if err != nil {
				log.Fatal(err)
			}
			f.Close()

			errDel := os.Remove(tempFile) // Remove the temp file, if there's an error we failed
			if errDel != nil {
				log.Fatal(errDel)
			}

			tempListing, err := os.ReadDir(tempPath) // Pull the listing for the subdirectory
			if err != nil {
				log.Fatal(err)
			}

			for _, subFile := range tempListing { // Check for the MAPPlaybackFile.dat which would denote we're in the session folder and that's a recording in the first folder
				if subFile.Name() == "MAPPlaybackFile.dat" {
					isThereRecording = true
					break // We found the file, that's enough
				}
			}
			break // We found a directory... either the file was there or not
		}
	}
	if !isThereRecording {
		log.Fatal("No MAPPlayback.dat file found!")
	}

}

func main() {

	// Sanity check - check the first child directory to see if we're in a recording session. Also check write/delete.
	// It will force os.Exit() if operations don't work.
	sanityCheck()

	// Go "while True" loop. Will exit when menuSelection is invalid
	for {
		menuSelection := mainMenu() // Show main menu to select function

		if menuSelection == 1 {
			prepRecordings() // User selected to prep recordings, run function
		} else if menuSelection == 2 {
			prepAnalysis() // User selected to prep analysis, run function
		} else {
			return // Leave the loop and function, we're done
		}
	}
}
