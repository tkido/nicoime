package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func publishHTML(ts []Tran) {
	const path = `nicoime_latest.html`
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	dt := time.Now().Format("2006年01月02日 15:04")
	const tmpl = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
<title>ニコニコ大百科IME辞書最新版データ</title>
</head>
<body>
現在%s版。<br />
登録単語数 %dです。<br />
<a href='http://tkido.com/data/nicoime.zip'>nicoime.zipをダウンロードする。</a>
</body>
</html>`
	html := fmt.Sprintf(tmpl, dt, len(ts))
	f.WriteString(html)

}

func publishAtok(ts []Tran) {
	const path = `nicoime_atok_utf8.txt`
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	dt := time.Now().Format("06/01/02 15:04")
	const headerTmpl = `!!ATOK_TANGO_TEXT_HEADER_1
!!一覧出力
!!対象辞書;nicoime.dic
!!単語種類;登録単語(*) 自動登録単語($)
!!読み範囲;(読みの先頭) → (読みの最終)
!!出力日時;%s

`
	header := fmt.Sprintf(headerTmpl, dt)
	f.WriteString(header)

	const format = "%s\t%s\t固有一般*\n"
	for _, t := range ts {
		fmt.Fprintf(f, format, strings.Replace(t.Read, `ゔ`, `う゛`, -1), t.Word)
	}
	return
}

func publishMs(ts []Tran) {
	const path = `nicoime_msime_utf8.txt`
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	dt := time.Now().Format("2006年01月02日 15:04")
	const headerTmpl = `!Microsoft IME Dictionary Tool 98
!Version:
!Format:WORDLIST
!User Dictionary Name: nicoime.dic
!Output File Name: nicoime_msime.txt
!DateTime: %s

`
	header := fmt.Sprintf(headerTmpl, dt)
	f.WriteString(header)

	const format = "%s\t%s\t固有名詞\n"
	for _, t := range ts {
		fmt.Fprintf(f, format, strings.Replace(t.Read, `ゔ`, `ヴ`, -1), t.Word)
	}
	return
}

func publish(ts []Tran) {
	publishHTML(ts)
	publishAtok(ts)
	publishMs(ts)
}
