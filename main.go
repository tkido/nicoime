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

	// fmt.Println(len(raws))
	_, err = process(raws)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(len(ts))
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
	if isNg(r) {
		return r, false
	}
	return convert(r), true
}

func convert(r Tran) Tran {
	// 「シリーズ」の除去
	if strings.HasSuffix(r.Word, `シリーズ`) {
		r.Word = r.Word[:len(r.Word)-12]
		if strings.HasSuffix(r.Read, `シリーズ`) {
			r.Read = r.Read[:len(r.Read)-12]
		}
	}
	// 「有限会社」等の除去
	word := r.Word
	word = strings.Replace(word, `有限会社`, ``, 1)
	word = strings.Replace(word, `(有)`, ``, 1)
	word = strings.Replace(word, `（有）`, ``, 1)
	if word != r.Word {
		r.Read = strings.Replace(r.Read, `ユウゲンガイシャ`, ``, 1)
	}
	r.Word = word
	// 「株式会社」等の除去
	word = strings.Replace(word, `株式会社`, ``, 1)
	word = strings.Replace(word, `(株)`, ``, 1)
	word = strings.Replace(word, `（株）`, ``, 1)
	if word != r.Word {
		switch {
		case strings.Contains(r.Read, `カブシキガイシャ`):
			r.Read = strings.Replace(r.Read, `カブシキガイシャ`, ``, 1)
		case strings.HasPrefix(r.Read, `マエカブ`):
			r.Read = strings.Replace(r.Read, `マエカブ`, ``, 1)
		case strings.HasSuffix(r.Read, `カッコカブ`):
			r.Read = strings.Replace(r.Read, `カッコカブ`, ``, 1)
		case strings.HasSuffix(r.Read, `カブ`):
			r.Read = strings.Replace(r.Read, `カブ`, ``, 1)
		}
	}
	r.Word = word
	// 「カッコ……カッコトジル」等の除去
	if reParenKana.MatchString(r.Read) {
		if !reParenKanaNg.MatchString(r.Read) {
			fmt.Println(r)
			r.Read = reParenKana.ReplaceAllString(r.Read, ``)
			r.Word = reParen.ReplaceAllString(r.Word, ``)
			fmt.Println(r)
		}
	}
	return r
}

var (
	reParen       = regexp.MustCompile(`(\(|（).*?(\)|）)$`)
	reParenKana   = regexp.MustCompile(`カッコ(.*?)(カッコ(トジル?)?)?$`)
	reParenKanaNg = regexp.MustCompile(`カッコ(イイ|ワルイ|ク|ヨス)(.*?)(カッコ(トジル?)?)?$`)
)

// type Pair struct {
// 	From, To string
// }

// var (
// 	pairsForWord = []Pair{
// 		{"シリーズ", ""},
// 	}
// )

// func replaceWordString(s string) string {
// 	return ""
// }
