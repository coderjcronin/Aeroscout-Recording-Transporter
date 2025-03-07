package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
	cleanRecordings

Iterate through the recording subdirectories, removing RadioMaps and Maps, preserving the first iteration for the ALE
*/
func cleanRecordings() {
	firstDirectory := true

	initialPath := filepath.Base(".")

	entries, err := os.ReadDir(initialPath)
	if err != nil {
		log.Panic(err) // if there's an error, log and exit (sanity check should of caught this but admin may have changed things on us during execution)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if firstDirectory { // Skip the first directory.
				log.Println("Skipping first directory, " + filepath.Join(initialPath, entry.Name()))
				firstDirectory = false
				continue
			}
			tempRadioMaps := filepath.Join(initialPath, entry.Name(), "RadioMaps")
			tempMaps := filepath.Join(initialPath, entry.Name(), "Maps")

			errRadioMaps := os.RemoveAll(tempRadioMaps)
			if errRadioMaps != nil {
				log.Panic(err) // Something didn't work, probably permissions or file in use
				return
			} else {
				log.Println("Deleted: " + tempRadioMaps) // Log and move on
			}

			errMaps := os.RemoveAll(tempMaps)
			if errMaps != nil {
				log.Panic(err) // Something didn't work, probably permissions or file in use
				return
			} else {
				log.Println("Deleted: " + tempMaps) // Log and move on
			}

		}
	}
	log.Println("Session cleaning has been completed.")
	os.Exit(0)
}

/*
	prepTransfer

Send contents to session folder (excluding ART) to a zip archive outside the session folder
*/
func prepTransfer() {
	// Get session name (it's the parent folder)
	fullDir, _ := os.Getwd()
	parentDir := filepath.Base(fullDir)
	zipPath := filepath.Join("..", parentDir+".zip") // Create zip file outside parent folder (help avoid having to skip it)

	file, err := os.Create(zipPath) // Create the zip file
	if err != nil {
		log.Panic(err) // Panic on error and exit
		return
	}

	defer file.Close() // Defer file close... we don't need it closed now but we will

	w := zip.NewWriter(file) // Setup a new writer into the zip file we made
	defer w.Close()          // Defer file close

	// local function for adding files to the zip file as they're walked by filepath.Walk()
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // We shouldn't start with an error...
		}
		if info.IsDir() {
			return nil // We don't need empty directories
		}

		pathInZip := strings.TrimLeft(path, fmt.Sprintf("..%c", os.PathSeparator)) // Correct the path for the zip folder... it gets messy doing ../ in a zip

		file, err := os.Open(path)
		if err != nil {
			return err // If we can't open the file to add to the zip, error out
		}
		defer file.Close()

		if filepath.Base(file.Name()) == "ART.exe" || filepath.Base(file.Name()) == "log.log" { // We don't want to zip this program or it's log
			return nil
		}

		f, err := w.Create(pathInZip)
		if err != nil {
			return err // If we can't add to the zip, error out
		}

		_, err = io.Copy(f, file) // Finally, copy the file into the zip archive
		if err != nil {
			return err // If it didn't work, error out
		} else {
			log.Println("Adding: " + pathInZip) // Otherwise, log the add and move on
		}

		return nil
	}
	err = filepath.Walk(filepath.Join("..", parentDir), walker) // Walk the directory path starting at parent directory, calling the local walker() function to add to zip
	if err != nil {
		log.Panic(err)
		return
	}

}

/*
	mainMenu - Returns int based on selection

Simply displays the utility main menu and returns an int indicating the user's selection.
1 - Process recordings for transport
2 - Process recordings for analysis
Anything else - Exit
*/
func mainMenu() int {

	var returnValue int

	fmt.Print("Please select an operation mode:\n\t1) Clean recording session\n\t2) Prepare session for transfer or storage\n\nPlease type the number followed by the ENTER key (invalid response will exit): ")
	_, err := fmt.Scanf("%d", &returnValue)

	if err != nil {
		log.Panic(err)
		return 0
	}

	return returnValue
}

/*
	sanityCheck() - returns boolean based on environment values

sanityCheck will check that:
1. We're in a recording session folder (there's a child folder with MAPPlaybackFile.dat in the first directory we walk to)
2. That we can create a file
3. That we can delete the file we created

If we fail at any point, log.Fatal() is raised which will cause an os.Exit() and terminate the program
*/
func sanityCheck() {

	isThereRecording := false
	baseDirectory := filepath.Base(".")

	entries, err := os.ReadDir(baseDirectory) // Pull the file listing for the local directory

	// If there was an error, log (well, print) it and return false
	if err != nil {
		log.Fatal(err)
		return
	}

	// Iterate through the listing for a directory
	for _, entry := range entries {
		if entry.IsDir() { // Found a directory, let's do a write and delete check, then look for our file!
			tempPath := filepath.Join(baseDirectory, entry.Name())
			tempFile := filepath.Join(tempPath, "ART.tst")

			f, err := os.Create(tempFile) // Create the temp file, if there's an error we failed
			if err != nil {
				log.Fatal(err)
				return
			}
			f.Close()

			errDel := os.Remove(tempFile) // Remove the temp file, if there's an error we failed
			if errDel != nil {
				log.Fatal(errDel)
				return
			}

			tempListing, err := os.ReadDir(tempPath) // Pull the listing for the subdirectory
			if err != nil {
				log.Fatal(err)
				return
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
		return
	}

}

/*
	main

Main call
*/
func main() {

	// Setup log to file
	logFile := filepath.Join(".", "log.log")

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//defer close when we're done
	defer f.Close()

	//set logger to use the f (logfile)
	log.SetOutput(f)

	//log and move on
	log.Println("Logging set to log file: " + logFile)

	// Sanity check - check the first child directory to see if we're in a recording session. Also check write/delete.
	// It will force os.Exit() if operations don't work.
	sanityCheck()

	// Main Menu selection
	menuSelection := mainMenu() // Show main menu to select function
	if menuSelection == 1 {
		cleanRecordings() // User selected to clean recordings, run function
	} else if menuSelection == 2 {
		prepTransfer() // User selected to prep for transfer or storage, run function
	}
}
