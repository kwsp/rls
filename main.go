package main

import (
	"fmt"
	"io/ioutil"
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

func tree(name string, prefix string) error {
	node, err := os.Lstat(name)
	if err != nil {
		return err
	}

	// Check if file is a symlink
	// Print name
	var sym_name string
	if node.Mode()&os.ModeSymlink != 0 {
		// Get symlink path
		sym_name, err = os.Readlink(name)
		if err != nil {
			return err
		}
		// Print current name and prefix
		fmt.Println(prefix + node.Name() + " -> " + sym_name)
	} else {
		// Print current name and prefix
		fmt.Println(prefix + node.Name())
	}

	// if node is a file, increment f_counter and return
	if !node.IsDir() {
		f_count++
		return nil
	}

	// node is a directory
	d_count++

	// Adjust the prefix for subdirectories
	if len(prefix) == 0 {
		prefix = TBAR
	} else {
		if strings.HasSuffix(prefix, LBAR) {
			prefix = prefix[:len(prefix)-10] + SPCE + TBAR
		} else {
			prefix = prefix[:len(prefix)-10] + VBAR + TBAR
		}
	}

	// Read list of directory entries
	dir_files, err := ioutil.ReadDir(name)
	if err != nil {
		println("Oops\n")
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
			prefix = prefix[:len(prefix)-10] + LBAR
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

	fmt.Printf("\n%d directories, %d files.\n", d_count, f_count)
}
