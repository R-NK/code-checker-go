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
	"github.com/yuin/charsetutil"
)

type opts struct {
	Status     bool     `short:"s" long:"status" description:"show files encoding and EOL"`
	Replace    string   `short:"r" long:"replace" description:"replace EOL with a specified argument. e.g. r=CRLF, r=LF" optional:"true"`
	TargetDir  string   `short:"t" long:"target" description:"target dir" default:"."`
	TargetExts []string `short:"e" long:"exts" description:"target file extensions" default:"*"`
}

func main() {
	var options opts
	parser := flags.NewParser(&options, flags.Default)
	_, err := parser.Parse()

	// コマンドライン引数が与えられない場合helpを表示
	if options.Status || options.Replace != "" {
		run(options)
	} else {
		if err == nil {
			parser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func run(options opts) {
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
			// relative file path
			fmt.Println(file)
			// file EOL
			fmt.Println(encoding)
			if strings.Contains(str, "\r\n") {
				fmt.Println("EOL: CRLF")
			} else {
				fmt.Println("EOL: LF")
			}
			fmt.Println()
		}

		if options.Replace != "" {
			var newEol string
			if options.Replace == "CRLF" {
				newEol = "\r\n"
			} else if options.Replace == "LF" {
				newEol = "\n"
			} else {
				log.Fatalln("ivalid argument")
			}
			rep := regexp.MustCompile(`\r\n|\r|\n`)
			replaced := rep.ReplaceAllString(str, newEol)
			ioutil.WriteFile("output/"+filepath.Base(file)+"_new", []byte(replaced), os.ModePerm)
			fmt.Println(file + "\n" + "replaced eol to " + options.Replace + "\n")
		}
	}
}
