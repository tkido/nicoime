package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	lines := readLines("initials.txt")
	fmt.Printf("lines: %v\n", lines)
}

func readLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File \"%s\" could not read: %v\n", path, err)
		os.Exit(1)
	}
	defer f.Close()

	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if serr := scanner.Err(); serr != nil {
		fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", path, serr)
	}

	return lines
}
