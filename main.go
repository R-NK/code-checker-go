package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/saintfish/chardet"
)

func main() {
	flag.Parse()

	content, err = ioutill.ReadFile(flag.Arg(0))

	var fp *os.File
	var err error

	fp, err = os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
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
