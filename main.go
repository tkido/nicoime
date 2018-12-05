package main

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// ZB is zero byte struct
type ZB struct{}

const (
	unit     = 50
	rawPath  = "data/raw.txt"
	tranPath = "data/tran.txt"
)

var (
	raws  = Trans{}
	trans = Trans{}
)

func main() {
	err := Load(rawPath, &raws)
	if err != nil {
		log.Fatal(err)
	}
	// cs, err := getCapitals()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = getRawData(cs)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = Save(rawPath, raws)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println(len(raws))
	ts, err := process(raws)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(ts))
}

func process(raws Trans) (ts Trans, err error) {
	ts = make(Trans, 0)

	wg := sync.WaitGroup{}
	q := make(chan ZB, runtime.NumCPU())
	ch := make(chan Tran)
	for _, r := range raws {
		wg.Add(1)
		go func(r Tran) {
			q <- ZB{}
			defer func() { <-q; wg.Done() }()
			if t, ok := filter(r); ok {
				ch <- t
			}
		}(r)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for t := range ch {
		ts[t.Word] = t
	}
	return
}

func filter(r Tran) (t Tran, ok bool) {
	if isNgRegex(r) || isNgLength(r) || isNgSuffix(r) || isNgRead(r) {
		return r, false
	}
	return r, true
}

var (
	reDate   = regexp.MustCompile(`^[平成昭和一二三四五六七八九十0-9\.:年月日]+$`)
	reDaikai = regexp.MustCompile(`第.*?回`)
)

func isNgRegex(r Tran) bool {
	if reDaikai.MatchString(r.Word) || reDate.MatchString(r.Word) {
		fmt.Println(r.Word)
		return true
	}
	return false
}

func isNgLength(r Tran) bool {
	if len([]rune(r.Read)) <= 3 {
		// fmt.Println(r.Word)
		return true
	}
	return false
}

var ngSuffixes = []string{
	`リンク`,
	`リンク集`,
	`一覧`,
	`のサムネ画像集`,
	`のお絵カキコ`,
}

func isNgSuffix(r Tran) bool {
	for _, suf := range ngSuffixes {
		if strings.HasSuffix(r.Word, suf) {
			// fmt.Println(r.Word)
			return true
		}
	}
	return false
}

var ngReads = []string{
	`ダイジョウブカ`,
	`ミエタ`,
	`イマ`,
	`アリガトウゴザイマス`,
	`オネガイシマス`,
	`オーケー`,
	`オキマシタ`,
	`ムリ`,
	`ムリデス`,
	`オダイジニ`,
	`ツウジョウエイギョウ`,
	`キヅカナカッタ`,
	`オワッタ`,
	`オツカレサマデス`,
	`イインジャネ`,
	`ホア`,
	`マウ`,
	`ヤッタカ`,
	`ジャナイ`,
	`イラナイ`,
	`チガウ`,
	`オモイダシタ`,
	`ドウシヨウモナイ`,
}

func isNgRead(r Tran) bool {
	for _, read := range ngReads {
		if r.Read == read {
			// fmt.Println(r.Word)
			return true
		}
	}
	return false
}
