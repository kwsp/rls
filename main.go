package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var fCount int = 0
var dCount int = -1

var _Space string = "    "
var _Vbar string = "│   "
var _Tbar string = "├── "
var _Lbar string = "└── "

/*
Check if the entry name is valid
*/
func toScan(name string) bool {
	if name[0] == '.' {
		return false
	}
	return true
}

func tree(name string, prefix string) error {
	node, err := os.Lstat(name)

	if err != nil {
		return err
	}

	// Check if file is a symlink
	// Print name
	var symName string
	if node.Mode()&os.ModeSymlink != 0 {
		// Get symlink path
		symName, err = os.Readlink(name)
		if err != nil {
			return err
		}

		// Print current name and prefix + symlink target
		fmt.Println(prefix + node.Name() + " -> " + symName)

		if symName[0] != '/' {
			symName = path.Join(filepath.Dir(name), symName)
		}

		node, err = os.Lstat(symName)
		if err != nil && os.IsNotExist(err) {
			fCount++
		} else {
			if node.IsDir() {
				dCount++
			} else {
				fCount++
			}
		}
		return nil

	}

	// Not symlink
	// Print current name and prefix
	fmt.Println(prefix + node.Name())

	// if node is a file, increment fCounter and return
	if !node.IsDir() {
		fCount++
		return nil
	}

	// node is a directory
	dCount++

	// Adjust the prefix for subdirectories
	if len(prefix) == 0 {
		prefix = _Tbar
	} else {
		if strings.HasSuffix(prefix, _Lbar) {
			prefix = prefix[:len(prefix)-10] + _Space + _Tbar
		} else {
			prefix = prefix[:len(prefix)-10] + _Vbar + _Tbar
		}
	}

	// Read list of directory entries
	dirFiles, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}

	// Purge dotfiles and directorys from the list in place
	n := 0
	for _, dirFile := range dirFiles {
		if toScan(dirFile.Name()) {
			dirFiles[n] = dirFile
			n++
		}
	}
	dirFiles = dirFiles[:n]

	// Traverse the files in the directory
	nFiles := len(dirFiles) - 1
	for i, dirFile := range dirFiles {

		// Change prefix for last entry
		if i == nFiles {
			prefix = prefix[:len(prefix)-10] + _Lbar
		}

		// Recursively call tree on each valid entry
		err = tree(path.Join(name, dirFile.Name()), prefix)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var dir string

	if len(os.Args) == 1 {
		dir = "."
	} else {
		dir = os.Args[1]
	}

	// Traverse
	err := tree(dir, "")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n%d directories, %d files\n", dCount, fCount)
}
