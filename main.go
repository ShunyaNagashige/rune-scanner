package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode/utf8"
)

type RuneScanner struct{
	r io.Reader

	//(参考)あくまでも、我々が扱っているのは、UTF-8により符号化された結果のバイト列。
	//なので、一バイト文字の場合、当然16文字までおさまる。
	//×訂正↑(参考)Unicodeのcode pointは4byte,最大1114112文字分
	//16byte→コードポイント４個分。今回の"Hello,World"は11文字であり、収まらない。
	//なので、のちに、io.MultiReaderにより
	//「bufで余った分」+「まだr(つまり*strings.Reader)に残っている分」
	//してくれる。つまり、連結(concat)してくれる
	buf [16]byte
}

//RuneScannerのポインタを返す
//ポインタを扱うことで、RuneScanner.bufへの破壊的な(関数をでてもその影響が残るような)操作が可能
func NewRuneScanner(r io.Reader)*RuneScanner{
	return &RuneScanner{r:r}
}

//(参考)runeの定義を調べてみる
//type rune = int32
func (s *RuneScanner)Scan()(rune,error){
	//buf[:]はbuf[0:len(buf)]と同じ
	//与えられたバイトスライスを先頭から埋めていく(rが*string.Readerであり、bufを"Hello,World"で埋める)
	//n:埋まったバイト数、err:埋める過程で発生したエラー
	n,err:=s.r.Read(s.buf[:])
	fmt.Printf("%s\n",s.buf)
	if err!=nil{
		return 0,err
	}

	//utf-8でエンコードされたバイト列から、
	//コードポイント(ルーン)をデコードする→「１コードポイントずつ」読み込む、を満たす
	//r:ルーン、size:バイトでの長さ
	r,size:=utf8.DecodeRune(s.buf[:n])
	if r==utf8.RuneError{
		return 0,errors.New("RuneError")
	}

	//s.bufで余った分と、まだr(つまり*strings.NewReader)に残っている分を連結
	//io.MultiReaderは、複数のio.Readerを一つにまとめてくれる(concat)
	//例えば複数ファイルを読み込むときに、
	//「複数であることを意識しなくていい」ようにしてくれるらしい
	s.r=io.MultiReader(bytes.NewReader(s.buf[size:n]))
	return r,nil
}

func main(){
	//io.Readerを実装しているものならなんでも引数として渡して良し
	s:=NewRuneScanner(strings.NewReader("Hello,WorldNaga"))
	for{
		r,err:=s.Scan()

		//エラー処理を、まとめて、一箇所に書いている
		//多分、エラー処理を一度に同時にやれ、というわけではないと思う
		if err==io.EOF{
			break
		}

		if err!=nil{
			log.Fatal(err)
		}

		fmt.Printf("%c\n",r)
	}
}