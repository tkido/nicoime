package main

import (
	"runtime"
	"sort"
	"strings"
	"sync"
)

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
