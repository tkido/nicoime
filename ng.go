package main

import (
	"fmt"
	"regexp"
	"strings"
)

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
