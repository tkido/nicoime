package main

import (
	"bufio"
	"fmt"
  "io"
	"log"
	"net/http"
	"os"
)

func main() {
	lines := readLines("initials.txt")
	fmt.Printf("lines: %v\n", lines)

  url := fmt.Sprintf("http://dic.nicovideo.jp/m/yp/a/%s/1-", lines[0])
  fmt.Println(url)
  download(url)
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

func download(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
    f, err := os.Create("html/a.html")
    defer f.Close()
    if err != nil {
      log.Fatal(err)
    }
    _, err1 := io.Copy(f, res.Body)
    if err1 != nil {
      log.Fatal(err1)
    }
  }
}
