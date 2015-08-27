// main.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	UTF8      = 65001
	SHIFT_JIS = 932
)

func main() {
	var text string
	text, ok := readPipe()
	if ok ==false {
		text, ok = readFileByArg()
		if ok == false {
			os.Exit(1)
		}
	}

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


