package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var startFolder = "testdata"
var isRecursive = true
var oldString = "_rep_"
var newString = "_replaced_"
var targetFolder = "target"

// Main function of application printing filenames
func main() {
	fmt.Println("Starting application...")
	err := os.Mkdir(targetFolder, os.ModeDir)
	if err != nil {
		fmt.Println(err)
	}
	err = printFiles(startFolder, isRecursive)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Ending application")
}

// function to print files in a given folder, inclusive subfolders if recursive
// behavior is enabled
func printFiles(folder string, isRecursive bool) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, f := range files {
		// Full file or foldername
		fullname := folder + "/" + f.Name()
		newName := strings.Replace(strings.Replace(fullname, oldString, newString, -1), startFolder+"/", "", 1)
		fmt.Println(newName)
		if f.IsDir() && isRecursive {
			err := os.Mkdir(targetFolder+"/"+newName, os.ModeDir)
			if err != nil {
				return err
			}
			err = printFiles(fullname, isRecursive)
			if err != nil {
				return err
			}
		} else {
			file, err := ioutil.ReadFile(fullname)
			if err != nil {
				return err
			}
			ioutil.WriteFile(targetFolder+"/"+newName, file, 0600)
		}
	}
	return nil
}
