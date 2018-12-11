package main

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func getRawData(cs []Capital) (err error) {
	wg := sync.WaitGroup{}
	q := make(chan ZB, runtime.NumCPU())
	ch := make(chan Tran)
	for _, c := range cs {
		for i := 1; i <= c.Count; i += unit {
			url := fmt.Sprintf(`https://dic.nicovideo.jp/m/yp/a/%s/%d-`, url.QueryEscape(c.Label), i)
			wg.Add(1)
			go func(url string) {
				q <- ZB{}
				defer func() { <-q; wg.Done() }()
				fmt.Printf("Downloading: %s\n", url)
				// err := download(url, ch)
				// if err != nil {
				// 	log.Println(err)
				// }
			}(url)
		}
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for r := range ch {
		raws[r.Word] = r
	}
	return
}

func download(url string, ch chan Tran) (err error) {
	resp, err := Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}
	lis := doc.Find(`#main > div.left-box > div > ul > ul > li`)
	lis.Each(func(i int, li *goquery.Selection) {
		word := li.Find(`a`).First().Text()

		lines := strings.Split(li.Text(), "\n")
		tail := lines[1][len(word)+12:]
		read := tail[:strings.Index(tail, ")")]

		redirect := strings.HasSuffix(lines[2], `(リダイレクト)`)

		ch <- Tran{word, read, "", redirect}
	})
	return
}
