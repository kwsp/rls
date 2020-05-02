package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var f_count int = 0
var d_count int = -1

var SPCE string = "    "
var VBAR string = "│   "
var TBAR string = "├── "
var LBAR string = "└── "

func toScan(name string) bool {
	if name[0] == '.' {
		return false
	}
	return true
}

func tree(name string, prefix []string) error {
	node, err := os.Stat(name)
	if err != nil {
		return err
	}

	// Print current name and prefix
	fmt.Println(strings.Join(prefix, "") + node.Name())

	// if node is a file, increment f_counter and return
	if !node.IsDir() {
		f_count++
		return nil
	}

	// node is a directory
	d_count++

	// Adjust the prefix for subdirectories
	if len(prefix) == 0 {
		prefix = []string{TBAR}
	} else {
		prefix[len(prefix)-1] = VBAR
		prefix = append(prefix, TBAR)
	}

	// Read its files
	dir, err := os.Open(name)
	if err != nil {
		return err
	}
	dir_files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// Traverse the files in the directory
	n_files := len(dir_files) - 1
	for i, dir_file := range dir_files {

		// Skip dotfiles and directories
		if !toScan(dir_file.Name()) {
			continue
		}

		// Change prefix for last entry
		if i == n_files {
			prefix = append(prefix[:len(prefix)-1], LBAR)
		}

		// Recursively call tree on each valid entry
		err = tree(path.Join(name, dir_file.Name()), prefix)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var dir string
	var prefix = []string{}

	if len(os.Args) == 1 {
		dir = "."
	} else {
		dir = os.Args[1]
	}

	// Traverse
	err := tree(dir, prefix)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n%d directories, %d files.\n", d_count, f_count)
}
