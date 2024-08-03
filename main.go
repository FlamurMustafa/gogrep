package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func getFiles(dirname string) ([]string, error) {
	var fileArray []string
	visitClo := func(filename string) fs.WalkDirFunc {
		return func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && d.Name() == ".git" {
				return fs.SkipDir
			}
			if !d.IsDir() {
				fileArray = append(fileArray, d.Name())
			}
			return nil
		}
	}

	e := filepath.WalkDir(dirname, visitClo(dirname))
	if e != nil {
		fmt.Print("Can't walk dir")
		return nil, e
	}
	return fileArray, nil
}

var tableMap map[rune]int

func createTableBadChar(str string) {
	tableMap = make(map[rune]int)

	for i, char := range str {
		tableMap[char] = i
	}
}

func match(pattern string, runes []rune) bool {
	m := len(pattern)
	n := len(runes)

	i := 0
	for i <= n-m {
		j := m - 1

		for j >= 0 && rune(pattern[j]) == runes[i+j] {
			j--
		}

		if j < 0 {
			return true
			i += m
		} else {
			charMisMatch := runes[i+j]
			shiftAmount := j - tableMap[charMisMatch]
			if shiftAmount < 1 {
				shiftAmount = 1
			}
			i = i + shiftAmount
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Incorrect Usage")
		fmt.Println("./gogrep \"string\" \"file\"")
		fmt.Println("OR ./gogrep \"string\" \"directory\" -r")
		os.Exit(-1)
	}

	var allFiles []string
	isRecursive := false
	substr := os.Args[1]
	fileName := os.Args[2]
	allFiles = append(allFiles, fileName)
	dirName := "."

	createTableBadChar(substr)

	for _, arg := range os.Args[1:] {
		if arg == "-r" {
			dirName = os.Args[2]
			isRecursive = true
		}
	}

	if isRecursive {
		allFiles, _ = getFiles(dirName)
	}

	lineNumber := 1
	for _, arg := range allFiles {
		file, err := os.Open(arg)
		if err != nil {
			fmt.Println("Can't open file\nAborting..")
			panic(err)
		}

		reader := bufio.NewReader(file)

		for {
			line, _, err := reader.ReadLine()
			lineNumber++
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				log.Fatal(err)
			}

			if match(substr, []rune(string(line))) {

				println(arg + ": " + strconv.Itoa(lineNumber))
			}
		}
		lineNumber = 0
	}
}
