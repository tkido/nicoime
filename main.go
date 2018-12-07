package main

import (
	"fmt"
	"log"
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
	// "シリーズ"の除去
	if strings.HasSuffix(r.Word, `シリーズ`) {
		fmt.Println(r)
		r.Word = r.Word[:len(r.Word)-12]
		if strings.HasSuffix(r.Read, `シリーズ`) {
			r.Read = r.Read[:len(r.Read)-12]
		}
		fmt.Println(r)
	}
	return r
}

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
