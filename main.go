package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/saintfish/chardet"
	"github.com/yuin/charsetutil"
)

var (
	status = *flag.Bool("status", false, "show files encoding and EOL")
	target = *flag.String("t", ".", "target dir")
)

func main() {
	flag.Parse()

	// exts, err := ioutil.ReadFile(flag.Arg(0))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	files := listFilesByExts(target, []string{"cpp", "h"})

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(file)
		encoding := detectCharEncode(content)
		fmt.Println(encoding)

		b, err := charsetutil.DecodeBytes(content, encoding)
		str := string(b)

		if strings.Contains(str, "\r\n") {
			fmt.Println("EOL: CRLF")
		} else {
			fmt.Println("EOL: LF")
		}
		fmt.Println()
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

func listFilesByExts(dir string, exts []string) []string {
	paths := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
			return err
		}
		for _, ext := range exts {
			// 拡張子が.から始まらない場合先頭に付与
			if !strings.Contains(ext, ".") {
				ext = "." + ext
			}
			if ext == filepath.Ext(path) {
				// fmt.Println(path)
				paths = append(paths, path)
			}
		}
		return nil
	})
	fmt.Println()
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
	}
	return paths
}

func replaceNewline(str string, newline string, oldline string) string {
	return strings.NewReplacer(oldline, newline).Replace(str)
}
