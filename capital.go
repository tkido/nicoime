package main

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// Capital is capital
type Capital struct {
	Label string
	Count int
}

func getCapitals() (cs []Capital, err error) {
	cs = make([]Capital, 0)
	resp, err := Get("http://dic.nicovideo.jp/m/a/a")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}
	trs := doc.Find(`#main > div.st-box > div > table:nth-child(3) > tbody > tr`)
	trs.Each(func(i int, tr *goquery.Selection) {
		tds := tr.Find(`td`)
		tds.Each(func(i int, td *goquery.Selection) {
			as := td.Find(`a`)
			label := as.First().Text()
			numS := as.Last().Text()
			num, err := strconv.Atoi(numS[1 : len(numS)-1])
			if err != nil {
				return
			}
			cs = append(cs, Capital{label, num})
		})
	})
	return
}
