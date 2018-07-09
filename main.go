package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/saintfish/chardet"
	"github.com/yuin/charsetutil"
)

type opts struct {
	Status     bool     `short:"s" long:"status" description:"show files encoding and EOL"`
	Replace    string   `short:"r" long:"replace" description:"replace EOL with a specified argument. e.g. r=CRLF, r=LF" optional:"true" default:"disable"`
	TargetDir  string   `short:"t" long:"target" description:"target dir" default:"."`
	TargetExts []string `short:"e" long:"exts" description:"target file extensions"`
}

func main() {
	var options opts
	_, err := flags.Parse(&options)
	if err != nil {
		fmt.Println("なんかおきたぞ")
	}

	files := listFilesByExts(options.TargetDir, options.TargetExts)
	if len(files) == 0 {
		fmt.Println("file not found.")
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		encoding := detectCharEncode(content)

		b, err := charsetutil.DecodeBytes(content, encoding)
		str := string(b)

		if options.Status {
			fmt.Println(file)
			fmt.Println(encoding)
			if strings.Contains(str, "\r\n") {
				fmt.Println("EOL: CRLF")
			} else {
				fmt.Println("EOL: LF")
			}
			fmt.Println()
		}

		if options.Replace != "disable" {
			var eol string
			if options.Replace == "CRLF" {
				eol = "\r\n"
			} else if options.Replace == "LF" {
				eol = "\n"
			} else {
				log.Fatalln("ivalid argument")
			}
			rep := regexp.MustCompile(`\r\n|\r|\n`)
			replaced := rep.ReplaceAllString(str, eol)
			ioutil.WriteFile("output/"+filepath.Base(file)+"_new", []byte(replaced), os.ModePerm)
			fmt.Println(file + "\n" + "replaced eol to " + options.Replace + "\n")
		}
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
