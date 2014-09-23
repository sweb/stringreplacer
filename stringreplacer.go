package main

import (
	"fmt"
	"io/ioutil"
)

var startFolder = "testdata"

// Main function of application printing filenames
func main() {
	fmt.Println("Starting application...")
	printFiles(startFolder)
	fmt.Println("Ending application")
}

// function to print files in a given folder, inclusive subfolders
func printFiles(folder string) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, f := range files {
		fmt.Println(folder + "/" + f.Name())
		if f.IsDir() {
			err := printFiles(folder + "/" + f.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
