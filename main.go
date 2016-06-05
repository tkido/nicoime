package main

import (
	"bufio"
	"fmt"
  "log"
	"os"
)

func main() {
	lines := readLines("initials.txt")
	fmt.Printf("lines: %v\n", lines)
}

func readLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
    log.Fatal(err)
	}
	defer f.Close()

	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
    log.Fatal(err)
	}

	return lines
}
