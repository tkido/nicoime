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
	r.Word = strings.Replace(r.Word, "シリーズ", "", 1)
	r.Read = strings.Replace(r.Read, "シリーズ", "", 1)
	return r, true
}

type Pair struct {
	From, To string
}

var (
	pairsForWord = []Pair{
		{"シリーズ", ""},
	}
)

func replaceWordString(s string) string {

}
