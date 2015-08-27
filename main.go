// main.go
package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
)

func main() {
	var text string
	var err error
	text, ok := readPipe()
	if ok != true {
		text, ok = readFileByArg()
		if ok == false {
			os.Exit(1)
		}
	}
	//encode
	text, err = transEnc(text)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	//dajarep
	s := Dajarep(text)
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i])
	}

}

func readPipe() (string, bool) {
	stats, _ := os.Stdin.Stat()

	if stats != nil && (stats.Mode()&os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err.Error())
			return "", false
		}

		return string(bytes), true
	} else {
		return "", false
	}
}
func readFileByArg() (string, bool) {

	if len(os.Args) < 2 {

		return "", false
	}
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}

	return string(content), true
}

//「Golangで文字コード判定」qiita.com/nobuhito/items/ff782f64e32f7ed95e43
func transEnc(text string) (asset string, err error) {
	body := []byte(text)
	var f []byte
	encodings := []string{"sjis", "utf-8"}
	for _, enc := range encodings {
		if enc != "" {
			ee, _ := charset.Lookup(enc)
			if ee == nil {
				continue
			}
			var buf bytes.Buffer
			ic := transform.NewWriter(&buf, ee.NewDecoder())
			_, err := ic.Write(body)
			if err != nil {
				continue
			}
			err = ic.Close()
			if err != nil {
				continue
			}
			f = buf.Bytes()
			break
		}
	}
	return string(f), nil
}
