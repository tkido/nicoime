package main

import (
	"log"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	nkf "github.com/creasty/go-nkf"
)

// ZB is zero byte struct
type ZB struct{}

const (
	unit    = 50
	rawPath = "raw.txt"
)

var (
	raws  = Trans{}
	trans = Trans{}
)

func main() {
	// err := Load(rawPath, &raws)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cs, err := getCapitals()
	if err != nil {
		log.Fatal(err)
	}
	err = getRawData(cs)
	if err != nil {
		log.Fatal(err)
	}
	// err = Save(rawPath, raws)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	ts, err := process(raws)
	if err != nil {
		log.Fatal(err)
	}
	publish(ts)
}

func process(raws Trans) (ts []Tran, err error) {
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
	tm := make(Trans, 0)
	for t := range ch {
		tm[t.Word] = t
	}

	count := make(map[string]int, 100000)
	for _, t := range tm {
		if note, ok := tm[t.Note]; ok {
			if strings.HasSuffix(t.Read, note.Read) {
				t.Read = t.Read[:len(t.Read)-len(note.Read)]
			}
		}
		count[t.Read]++
	}
	ts = make([]Tran, 0, 100000)
	for _, t := range tm {
		if !(count[t.Read] >= 2 && t.Redirect) {
			ts = append(ts, t)
		}
	}
	sort.SliceStable(ts, func(i, j int) bool {
		return ts[i].Read < ts[j].Read
	})
	return
}

func filter(r Tran) (t Tran, ok bool) {
	r = convert(r)
	if isNg(r) {
		return r, false
	}
	return r, true
}

func convert(r Tran) Tran {
	// 謎の記号の除去
	r.Word = strings.Replace(r.Word, `̣`, ``, -1)
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
	if strings.HasSuffix(r.Word, `)`) || strings.HasSuffix(r.Word, `）`) {
		// 括弧で終わるものが対象
		if strings.HasPrefix(r.Word, `(`) || strings.HasPrefix(r.Word, `（`) {
			// 先頭が括弧の場合はそのまま
		} else {
			if ss := reParen.FindStringSubmatch(r.Word); ss != nil {
				r.Word = strings.Replace(r.Word, ss[0], "", 1)
				r.Note = ss[2]
			}
			if reParenKana.MatchString(r.Read) {
				if !reParenKanaNg.MatchString(r.Read) {
					r.Read = reParenKana.ReplaceAllString(r.Read, ``)
				}
			}
		}
	}
	// ひらがなに変換
	hira, err := nkf.Convert(r.Read, "-m0 -W -w --hiragana")
	if err != nil {
		// log.Printf("can't convert %s by nkf", r.Read)
	} else {
		r.Read = hira
	}
	return r
}

var (
	reParen       = regexp.MustCompile(`(\(|（)(.*?)(\)|）)$`)
	reParenKana   = regexp.MustCompile(`カッコ(.*?)(カッコ(トジル?)?)?$`)
	reParenKanaNg = regexp.MustCompile(`カッコ(イイ|ワルイ|ク|ヨス)(.*?)(カッコ(トジル?)?)?$`)
)
