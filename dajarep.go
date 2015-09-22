// sharep.go
package main

import (
	"github.com/ikawaha/kagome"
	"math"
	"regexp"
	"strings"
)

var (
	token  *kagome.Tokenizer
)

//単語
type word struct {
	str   string
	kana  string
	wtype string
}

//文章
type sentence struct {
	str   string
	kana  string
	words []word
}

//駄洒落を返す
func Dajarep(text string) (dajares []string, debugStrs []string) {
	sentences := getSentences(text)
	for i := 0; i < len(sentences); i++ {
		if ok, kana := isDajare(sentences[i]); ok == true {
			dajares = append(dajares, sentences[i].str)
			debugStrs = append(debugStrs, kana)
		}
	}
	return dajares, debugStrs
}

//駄洒落かどうかを評価する。
func isDajare(sen sentence) (bool, string) {
	words := sen.words
	for i := 0; i < len(words); i++ {
		w := words[i]
		if w.wtype == "名詞" && len([]rune(w.kana)) > 1 {
			r_str := regexp.MustCompile(w.str)
			r_kana := regexp.MustCompile(fixWord(w.kana))
			hit_str := r_str.FindAllString(sen.str, -1)
			hit_kana := r_kana.FindAllString(sen.kana, -1)
			hit_kana2 := r_kana.FindAllString(fixSentence(sen.kana), -1)
			//ある単語における　原文の一致文字列数<フリガナでの一致文字列数　→　駄洒落の読みが存在
			if len(hit_str) < int(math.Max(float64(len(hit_kana)), float64(len(hit_kana2)))) {
				return true, w.kana
			}
		}
	}
	return false, ""
}

//置き換え可能な文字を考慮した正規表現を返す。
func fixWord(text string) string {
	text = strings.Replace(text, "ッ", "[ツッ]?", -1)
	text = strings.Replace(text, "ー", "[ー]?", -1)
	text = strings.Replace(text, "ァ", "[アァ]?", -1)
	text = strings.Replace(text, "ィ", "[イィ]?", -1)
	text = strings.Replace(text, "ゥ", "[ウゥ]?", -1)
	text = strings.Replace(text, "ェ", "[エェ]?", -1)
	text = strings.Replace(text, "ォ", "[オォ]?", -1)
	text = strings.Replace(text, "ャ", "[ヤャ]", -1)
	text = strings.Replace(text, "ュ", "[ユュ]", -1)
	text = strings.Replace(text, "ョ", "[ヨョ]", -1)
	return text
}

//本文から省略可能文字を消したパターンを返す。
func fixSentence(text string) string {
	text = strings.Replace(text, "ッ", "", -1)
	text = strings.Replace(text, "ー", "", -1)
	return text
}

//テキストからsentenceオブジェクトを作る。
func getSentences(text string) []sentence {
	var sentences []sentence
	t:= getTokenizer()

	text = strings.Replace(text, "。", "\n", -1)
	text = strings.Replace(text, ".", "\n", -1)
	text = strings.Replace(text, "?", "?\n", -1)
	text = strings.Replace(text, "!", "!\n", -1)
	text = strings.Replace(text, "？", "？\n", -1)
	text = strings.Replace(text, "！", "！\n", -1)
	senstr := strings.Split(text, "\n")

	for i := 0; i < len(senstr); i++ {
		tokens := t.Tokenize(senstr[i])
		var words []word
		var kana string
		for j := 0; j < len(tokens); j++ {
			tk := tokens[j]
			ft := tk.Features()
			if len(ft) > 7 {
				w := word{str: ft[6],
					kana:  ft[7],
					wtype: ft[0],
				}
				words = append(words, w)
				kana += ft[7]
			}
		}
		sentences = append(sentences,
			sentence{
				str:   senstr[i],
				words: words,
				kana:  kana,
			})
	}
	return sentences
}
//Tokenizerを取得
func getTokenizer() *kagome.Tokenizer{
	if token == nil {
		token = kagome.NewTokenizer()
	}
	return token
}