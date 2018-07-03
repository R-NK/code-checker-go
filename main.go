package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/saintfish/chardet"
)

func main() {
	flag.Parse()

	content, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(detectCharEncode(content))
}

func detectCharEncode(body []byte) string {
	det := chardet.NewTextDetector()
	result, err := det.DetectBest(body)
	if err != nil {
		return "Unknown"
	}
	return result.Charset
}
