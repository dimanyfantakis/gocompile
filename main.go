package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.LstdFlags|log.Lshortfile)
)

func main() {
	path := "foo.teeny"
	readfile := readFile(path)
	defer readfile.Close()

	program := ""
	scanner := bufio.NewScanner(readfile)

	for scanner.Scan() {
		line := scanner.Text()
		program += line + "\n"
	}

	t := tokenize(program)
	parse(t)
}

func readFile(path string) *os.File {
	readfile, err := os.Open(path)
	if err != nil {
		os.Exit(1)
	}
	return readfile
}
