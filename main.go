package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/saintfish/chardet"
	"github.com/yuin/charsetutil"
)

func main() {
	flag.Parse()

	content, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	encoding := detectCharEncode(content)
	fmt.Printf(encoding + "\n")

	b, err := charsetutil.DecodeBytes(content, encoding)
	str := string(b)
	fmt.Printf(str)

	if strings.Contains(str, "\r\n") {
		fmt.Printf("EOL: CRLF\n")
	} else {
		fmt.Printf("EOL: LF\n")
	}
}

func detectCharEncode(body []byte) string {
	det := chardet.NewTextDetector()
	result, err := det.DetectBest(body)
	if err != nil {
		return "Unknown"
	}
	return result.Charset
}
