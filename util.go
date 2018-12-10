package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Get is http.Get with retry
func Get(url string) (resp *http.Response, err error) {
	const max = 3
	for i := 1; ; i++ {
		resp, err = http.Get(url)
		if err != nil {
			return
		}
		sc := resp.StatusCode
		switch {
		case sc == 200:
			return
		case sc >= 500:
			if i < max {
				time.Sleep(time.Second) // retry after 1 second
				continue
			}
			return nil, fmt.Errorf("http error status %d", sc)
		case sc >= 400:
			return nil, fmt.Errorf("http error status %d", sc)
		}
	}
}

// Save data
func Save(path string, data Trans) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	for _, v := range data {
		f.WriteString(v.Word)
		f.WriteString("\t")
		f.WriteString(v.Read)
		f.WriteString("\t")
		f.WriteString(fmt.Sprintf("%v", v.Redirect))
		f.WriteString("\n")
	}
	return
}

// Load data
func Load(path string, data *Trans) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		ss := strings.Split(s.Text(), "\t")
		redirect, err := strconv.ParseBool(ss[2])
		if err != nil {
			return err
		}
		r := Tran{ss[0], ss[1], "", redirect}
		(*data)[r.Word] = r
	}
	return
}
