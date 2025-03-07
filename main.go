package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

/*
	prepRecordings

Iterate through the recording subdirectories, removing RadioMaps and Maps, preserving the first iteration for restoration
*/
func prepRecordings() {
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
	log.Println("Transport processing has been completed.")
	os.Exit(0)
}

/*
	prepAnalysis

Iterates through recording directorties, save the first directory path for RadioMaps and Maps, then create symlinks in all other subdirectories
*/
func prepAnalysis() {
	firstDirectory := true

	// Setup initial path and prep parent folders
	initialPath, _ := os.Getwd()
	parentRadioMaps, _ := os.Getwd()
	parentMaps, _ := os.Getwd()

	entries, err := os.ReadDir(initialPath)
	if err != nil {
		log.Panic(err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if firstDirectory { // First directory should contain pristine Maps and RadioMaps; set the parent variables and move on with the loop
				parentMaps = filepath.Join(initialPath, entry.Name(), "Maps")
				parentRadioMaps = filepath.Join(initialPath, entry.Name(), "RadioMaps")
				firstDirectory = false
				log.Println("Initial Radio Maps: " + parentRadioMaps)
				log.Println("Initial Maps: " + parentMaps)
				continue
			}

			tempMaps := filepath.Join(initialPath, entry.Name(), "Maps")
			tempRadioMaps := filepath.Join(initialPath, entry.Name(), "RadioMaps")

			errLinkMaps := os.Symlink(parentMaps, tempMaps)
			if errLinkMaps != nil {
				log.Panic(errLinkMaps) // Something didn't work creating the symlink
				return
			} else {
				log.Println("Created symlink: " + tempMaps) // It worked, log and move on
			}

			errLinkRadioMaps := os.Symlink(parentRadioMaps, tempRadioMaps)
			if errLinkRadioMaps != nil {
				log.Panic(errLinkRadioMaps) // Something didn't work creating the symlink
				return
			} else {
				log.Println("Created symlink: " + tempRadioMaps) // It worked, log and move on
			}

		}
	}
	log.Println("Analysis processing has been completed.")
	os.Exit(0)

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

	fmt.Print("Please select an operation mode:\n\t1) Prepare recordings for transport or storage\n\t2) Prepare recordings for analysis and reporting\n\nPlease type the number followed by the ENTER key (invalid response will exit): ")
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
1. We're in a recording session folder (there's a child folder with MAPPlaybackFile.dat)
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
	checkForAdmin

Return true if we can read C: (physicaldrive0), false if not so we can ask for UAC admin
*/
func checkForAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")

	return err == nil
}

/*
	runMeElevated

Ask to relaunch with UAC admin access
*/
func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	// Request UAC elevated execution if we're on Win32
	osType := runtime.GOOS
	if osType == "windows" {
		if !checkForAdmin() {
			log.Println("Not in UAC Admin, relaunching....")
			runMeElevated()
			os.Exit(0)
		}
	}

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

	// Go "while True" loop. Will exit when menuSelection is invalid
	for {
		menuSelection := mainMenu() // Show main menu to select function

		fmt.Println(menuSelection)

		if menuSelection == 1 {
			prepRecordings() // User selected to prep recordings, run function
		} else if menuSelection == 2 {
			prepAnalysis() // User selected to prep analysis, run function
		} else {
			return // Leave the loop and function, we're done
		}
	}
}
