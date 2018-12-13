package main

import (
	"regexp"
	"strings"

	nkf "github.com/creasty/go-nkf"
)

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
