// main.go
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/kurehajime/dajarep"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"runtime"
)

func main() {
	var text string
	var err error
	var encode string
	var debug bool
	var interactive bool
	var default_encoding string

	if runtime.GOOS == "windows" {
		default_encoding = "sjis"
	} else {
		default_encoding = "utf-8"
	}
	flag.StringVar(&encode, "e", default_encoding, "encoding")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.BoolVar(&interactive, "i", false, "interactive mode")

	flag.Parse()

	if interactive == true {
		fmt.Print("> ")
	} else if len(flag.Args()) == 0 {
		text, err = readPipe()
	} else if flag.Arg(0) == "-" {
		text, err = readStdin()
	} else {
		text, err = readFileByArg(flag.Arg(0))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if interactive == false {
		text, err := transEnc(text, encode)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		s, d := dajarep.Dajarep(text, debug)
		for i := 0; i < len(s); i++ {
			if !debug {
				fmt.Println(s[i])
			} else {
				fmt.Println(s[i] + "[" + d[i] + "]")
			}
		}
	} else {
		//interactive mode
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			if s.Err() != nil {
				fmt.Fprintln(os.Stderr, s.Err())
				os.Exit(1)
			}
			if s.Text() == "" {
				break
			}
			text := s.Text()
			_, d := dajarep.Dajarep(text, debug)
			if len(d) > 0 {
				for i := 0; i < len(d); i++ {
					fmt.Println("-> " + d[i])
				}
			} else {
				fmt.Println("")
			}
			fmt.Print("> ")
		}
	}
}

func readPipe() (string, error) {
	stats, _ := os.Stdin.Stat()
	if stats != nil && (stats.Mode()&os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	} else {
		return "", nil
	}
}
func readStdin() (string, error) {
	var text string
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		text += s.Text() + "\n"
	}
	if s.Err() != nil {
		return "", s.Err()
	}
	return text, nil
}

func readFileByArg(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
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
