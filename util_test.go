package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestGet(t *testing.T) {
	resp, err := Get("https://dic.nicovideo.jp/m/yp/a/%E3%82%AF/451-")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		t.Error(err)
	}
	lis := doc.Find(`#main > div.left-box > div > ul > ul > li`)
	lis.Each(func(i int, li *goquery.Selection) {
		word := li.Find(`a`).First().Text()

		lines := strings.Split(li.Text(), "\n")
		// for _, line := range lines {
		// 	fmt.Printf("【%s】\n", line)
		// }

		tail := lines[1][len(word)+12:]
		read := tail[:strings.Index(tail, ")")]

		redirect := strings.HasSuffix(lines[2], `(リダイレクト)`)

		tran := Tran{word, read, redirect}
		fmt.Println(tran)
	})
}
