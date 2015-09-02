// main.go
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
)

func main() {
	var text string
	var err error
	var ok bool = false
	var encode string

	flag.StringVar(&encode, "e", "", "encoding")
	flag.Parse()

	if len(flag.Args()) == 0 {
		text, ok = readPipe()
	} else if flag.Arg(0) == "-" {
		text, ok = readStdin()
	} else {
		text, ok = readFileByArg(flag.Arg(0))
	}
	if ok == false {
		os.Exit(1)
	}

	//encode
	text, err = transEnc(text, encode)
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
			fmt.Fprintln(os.Stderr, err.Error())
			return "", false
		}
		return string(bytes), true
	} else {
		return "", false
	}
}
func readStdin() (string, bool) {
	var text string
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		text += s.Text() + "\n"
	}
	if s.Err() != nil {
		fmt.Fprintln(os.Stderr, s.Err().Error())
		return "", false
	}
	return text, true
}

func readFileByArg(path string) (string, bool) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return "", false
	}
	return string(content), true
}

//「Golangで文字コード判定」qiita.com/nobuhito/items/ff782f64e32f7ed95e43
func transEnc(text string, encode string) (string, error) {
	body := []byte(text)
	var f []byte

	encodings := []string{"sjis", "utf-8"}
	if encode != "" {
		encodings = append([]string{encode}, encodings...)
	}
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
