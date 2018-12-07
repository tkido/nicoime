package main

import (
	"regexp"
	"strings"
)

func isNg(r Tran) (ok bool) {
	if isNgLength(r) || isNgSuffix(r) || isNgRead(r) || isNgSubStr(r) || isNgRegex(r) {
		return true
	}
	return false
}

var ngSubStrs = []string{
	`生主`,
	`生放送主`,
}

func isNgSubStr(r Tran) bool {
	for _, sub := range ngSubStrs {
		if strings.Contains(r.Word, sub) {
			// fmt.Println(r.Word)
			return true
		}
	}
	return false
}

var (
	reDate   = regexp.MustCompile(`^[平成昭和一二三四五六七八九十0-9\.:年月日]+$`)
	reDaikai = regexp.MustCompile(`第.*?回`)
)

func isNgRegex(r Tran) bool {
	if reDaikai.MatchString(r.Word) || reDate.MatchString(r.Word) {
		// fmt.Println(r.Word)
		return true
	}
	return false
}

// 読みが3文字以下の項目を排除
func isNgLength(r Tran) bool {
	if len([]rune(r.Read)) <= 2 {
		// fmt.Println(r.Word)
		return true
	}
	return false
}

// IME辞書というより事典っぽい項目の排除
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

// 読み遊びの排除
var ngReads = []string{
	`アリガトウゴザイマス`,
	`イインジャネ`,
	`イマ`,
	`イラナイ`,
	`オーケー`,
	`オモイダシタ`,
	`オキマシタ`,
	`オダイジニ`,
	`オツカレサマデス`,
	`オネガイシマス`,
	`オワッタ`,
	`ガンメンキジョウ`,
	`キヅカナカッタ`,
	`ジャナイ`,
	`ダイジョウブカ`,
	`チガウ`,
	`ツウジョウエイギョウ`,
	`ドウシヨウモナイ`,
	`ホア`,
	`マウ`,
	`ヤッタカ`,
	`ミエタ`,
	`ムリ`,
	`ムリデス`,
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
