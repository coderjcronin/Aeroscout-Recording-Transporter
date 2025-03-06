package main

import (
	"bufio"
	"fmt"
	"os"
	// "log" We'll use this later
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

/*
	sanityCheck() - returns boolean based on environment values

sanityCheck will check that:
1. We're in a recording session folder (there's a child folder with MAPPlaybackFile)
__ I may need to disable this during development. __
2. That we can create a file
3. That we can delete the file we created

If all 3 are valid, then return true. Otherwise, either we need this utility moved or admin rights and let the user know.
*/
func sanityCheck() bool {

	isThereRecording := false

	// Block to check for the MAPPlaybackFile.dat
	//
	entries, err := os.ReadDir("./") // Pull the file listing for the local directory

	// If there was an error, log (well, print) it and return false
	if err != nil {
		fmt.Print(err)
		return false
	}

	// Iterate through the listing for a directory
	for _, file := range entries {
		fmt.Println(file.Name())
		if file.IsDir() {
			tempPath := "./" + file.Name()
			tempListing, err := os.ReadDir(tempPath)
			if err != nil {
				fmt.Print(err)
				return false
			}
			for _, subFile := range tempListing {
				fmt.Println(subFile.Name())
				if subFile.Name() == "MAPPlaybackFile.dat" {
					isThereRecording = true
					break // We found the file, that's enough
				}
			}
			//break // We found a directory... either the file was there or not
		}
	}
	if !isThereRecording {
		fmt.Println("No MAPPlayback.dat file found!")
		return false
	}
	// End the MAPPlayback block

	tempFileForTesting := "./artTest.dat"

	// Write check block
	f, err := os.Create(tempFileForTesting)
	if err != nil {
		fmt.Print(err)
		return false
	}
	f.Close()

	// Delete check block
	errDel := os.Remove(tempFileForTesting)
	if errDel != nil {
		fmt.Print(errDel)
		return false
	}

	return true
}

func main() {

	// Sanity check - check the first child directory to see if we're in a recording session.
	weGood := sanityCheck()

	if !weGood {
		fmt.Println("Exiting.")
		return
	}

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
