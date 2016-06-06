package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func main() {
	makeInitials()
	//lines := readLines("initials.txt")

	//url := fmt.Sprintf("http://dic.nicovideo.jp/m/yp/a/%s/1-", lines[0])
	//download(url)
}

func makeInitials() {
	res, err := http.Get("http://dic.nicovideo.jp/m/a/a")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("initials.txt")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(res.Body)
	re := regexp.MustCompile(`<a href="/m/yp/a/(.*?)/1-">\((\d*?)\)</a></td>`)
	for scanner.Scan() {
		ms := re.FindStringSubmatch(scanner.Text())
		if ms != nil {
			char, _ := url.QueryUnescape(ms[1])
			io.WriteString(f, fmt.Sprintf("%s,%s\n", char, ms[2]))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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
