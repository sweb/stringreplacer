package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var isRecursive = flag.Bool("isRecursive", true, "False if only the children of the start folder should be renamed")
var oldString = flag.String("old", "", "String that is going to be replaced")
var newString = flag.String("new", "", "String that is replacing the old string")
var targetFolder = flag.String("targetFolder", "target",
	"Folder containing results, relatively stored to the current path")
var ignoreWinRestr = flag.Bool("ignoreWinRestr", false,
	"The windows file system API only allows absolute filenames up to 260 characters. If this parameter is set to true the restriction is ignored")

var stringLine = strings.Repeat("-", 80)
var targetFolderRequired = true

// Main function of application printing filenames
func main() {
	flag.Parse()
	printSeparatorToCLI()
	fmt.Println("Starting application stringreplacer...")
	printSeparatorToCLI()

	fmt.Println("Parameters:")
	startFolder, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	*targetFolder = startFolder + `\` + *targetFolder
	fmt.Println("startFolder: " + startFolder)
	fmt.Println("replace " + *oldString + " with " + *newString)
	fmt.Println("targetFolder: " + *targetFolder)
	printSeparatorToCLI()
	if *oldString == "" {
		fmt.Println("Please provide correct parameters...")
		os.Exit(1)
	}
	pathPrefix := ""
	if *ignoreWinRestr {
		pathPrefix = `\\?\`
	}
	// Rename and save content of start folder
	err = renameAndSave(startFolder, *isRecursive, startFolder, pathPrefix)
	if err != nil {
		fmt.Println("An error occured:")
		fmt.Println(err)
	}
	printSeparatorToCLI()
	fmt.Println("Ending application")
}

// function to print files in a given folder, inclusive subfolders if recursive
// behavior is enabled
func renameAndSave(folder string, isRecursive bool, startFolder string,
	pathPrefix string) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}
	// Create target folder
	if targetFolderRequired {
		err = os.Mkdir(*targetFolder, os.ModeDir)
		if err != nil {
			return err
		}
		targetFolderRequired = false
	}

	for _, f := range files {
		// Full file or foldername
		fullname := folder + `\` + f.Name()
		newName := strings.Replace(strings.Replace(fullname, *oldString, *newString, -1), startFolder+`\`, "", 1)
		fmt.Println("  " + fullname + "\n    -->\n  " + newName)
		writableObjectName := pathPrefix + *targetFolder + `\` + newName
		if f.IsDir() && isRecursive {
			err = os.Mkdir(writableObjectName, os.ModeDir)
			if err != nil {
				printAdviceInCaseOfLongFilename(writableObjectName)
				return err
			}
			err = renameAndSave(fullname, isRecursive, startFolder, pathPrefix)
			if err != nil {
				return err
			}
		} else {
			file, err := ioutil.ReadFile(pathPrefix + fullname)
			if err != nil {
				printAdviceInCaseOfLongFilename(writableObjectName)
				return err
			}

			err = ioutil.WriteFile(writableObjectName, file, 0600)
			if err != nil {
				printAdviceInCaseOfLongFilename(writableObjectName)
				return err
			}
		}
	}
	return nil
}

func printSeparatorToCLI() {
	fmt.Println(stringLine)
}

func printAdviceInCaseOfLongFilename(filename string) {
	if !*ignoreWinRestr && len(filename) > 260 {
		printSeparatorToCLI()
		fmt.Println("*** This error was probably caused by windows file system restrictions.\n*** You can try the parameter ignoreWinRestr (see -help)")
		printSeparatorToCLI()
	}
}
