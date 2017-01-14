package dajarep

import (
	"math"
	"regexp"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

func init() {
	tokenizer.SysDic()
}

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

//Dajarep :駄洒落を返す
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
			rStr := regexp.MustCompile(w.str)
			rKana := regexp.MustCompile(fixWord(w.kana))
			hitStr := rStr.FindAllString(sen.str, -1)
			hitKana := rKana.FindAllString(sen.kana, -1)
			hitKana2 := rKana.FindAllString(fixSentence(sen.kana), -1)
			//ある単語における　原文の一致文字列数<フリガナでの一致文字列数　→　駄洒落の読みが存在
			if len(hitStr) < int(math.Max(float64(len(hitKana)), float64(len(hitKana2)))) {
				return true, w.kana
			}
		}
	}
	return false, ""
}

//置き換え可能な文字を考慮した正規表現を返す。
func fixWord(text string) string {
	text = strings.Replace(text, "ッ", "[ツッ]?", -1)
	text = strings.Replace(text, "ァ", "[アァ]?", -1)
	text = strings.Replace(text, "ィ", "[イィ]?", -1)
	text = strings.Replace(text, "ゥ", "[ウゥ]?", -1)
	text = strings.Replace(text, "ェ", "[エェ]?", -1)
	text = strings.Replace(text, "ォ", "[オォ]?", -1)
	text = strings.Replace(text, "ズ", "[ズヅ]", -1)
	text = strings.Replace(text, "ヅ", "[ズヅ]", -1)
	text = strings.Replace(text, "ヂ", "[ジヂ]", -1)
	text = strings.Replace(text, "ジ", "[ジヂ]", -1)
	re := regexp.MustCompile("([アカサタナハマヤラワャ])ー")
	text = re.ReplaceAllString(text, "$1[アァ]?")
	re = regexp.MustCompile("([イキシチニヒミリ])ー")
	text = re.ReplaceAllString(text, "$1[イィ]?")
	re = regexp.MustCompile("([ウクスツヌフムユルュ])ー")
	text = re.ReplaceAllString(text, "$1[ウゥ]?")
	re = regexp.MustCompile("([エケセテネへメレ])ー")
	text = re.ReplaceAllString(text, "$1[エェ]?")
	re = regexp.MustCompile("([オコソトノホモヨロヲョ])ー")
	text = re.ReplaceAllString(text, "$1[ウゥオォ]?")
	text = strings.Replace(text, "ャ", "[ヤャ]", -1)
	text = strings.Replace(text, "ュ", "[ユュ]", -1)
	text = strings.Replace(text, "ョ", "[ヨョ]", -1)
	text = strings.Replace(text, "ー", "[ー]?", -1)
	return text
}

//本文から省略可能文字を消したパターンを返す。
func fixSentence(text string) string {
	text = strings.Replace(text, "ッ", "", -1)
	text = strings.Replace(text, "ー", "", -1)
	text = strings.Replace(text, "、", "", -1)
	text = strings.Replace(text, ",", "", -1)
	text = strings.Replace(text, "　", "", -1)
	text = strings.Replace(text, " ", "", -1)
	return text
}

//テキストからsentenceオブジェクトを作る。
func getSentences(text string) []sentence {
	var sentences []sentence
	t := tokenizer.New()

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
