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
	OutputDir  string   `short:"o" long:"out" description:"output dir" default:"$OVERRIDE"`
}

// UNKNOWN represents file encoding cannot be detected
const UNKNOWN = "Unknown"

// OVERRIDE represents default OutputDIr parameter
const OVERRIDE = "$OVERRIDE"

func main() {
	var options opts
	parser := flags.NewParser(&options, flags.Default)
	_, err := parser.Parse()

	// show help when no arguments provided
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

	if options.OutputDir != OVERRIDE {
		if _, err := os.Stat(options.OutputDir); os.IsNotExist(err) {
			if err := os.Mkdir(options.OutputDir, 0777); err != nil {
				log.Fatalln(err)
			}
		}
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
			// convert UTF-8 to original encoding
			converted := charsetutil.MustEncode(replaced, encoding)

			var outPath string
			if options.OutputDir == OVERRIDE {
				outPath = file
			} else {
				outPath = options.OutputDir + "/" + filepath.Base(file)
			}

			ioutil.WriteFile(outPath, []byte(converted), os.ModePerm)
			fmt.Println(file + "\n" + "replaced eol to " + options.Replace + "\n")
		}

	}
}
