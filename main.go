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
	"strconv"
	"strings"
)

const pageUnit = 50

func main() {
	makeInitials()
	download(17)
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

func readInitials() ([]string, []int) {
	f, err := os.Open("initials.txt")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	initials := make([]string, 0, 85)
	limits := make([]int, 0, 85)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ss := strings.Split(scanner.Text(), ",")
		initials = append(initials, ss[0])
		i, _ := strconv.Atoi(ss[1])
		limits = append(limits, (i-1)/pageUnit)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return initials, limits
}

func download(hour int) {
	is, ls := readInitials()
	for i := hour * 4; i < hour*4+4; i++ {
		for j := 0; j <= ls[i]; j++ {
			num := j*pageUnit + 1
			url := fmt.Sprintf("http://dic.nicovideo.jp/m/yp/a/%s/%d-", url.QueryEscape(is[i]), num)
			path := fmt.Sprintf("data/%02d_%06d.html", i, num)
			res, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			f, err := os.Create(path)
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
}
