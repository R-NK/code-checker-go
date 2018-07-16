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
	Replace    string   `short:"r" long:"replace" description:"replace EOL with a specified argument. e.g. r=CRLF, r=LF"`
	TargetDir  string   `short:"t" long:"target" description:"target dir" default:"." default-mask:"current dir"`
	TargetExts []string `short:"e" long:"exts" description:"target file extensions" default:"*" default-mask:"all"`
	OutputDir  string   `short:"o" long:"out" description:"output dir" default:"$OVERRIDE" default-mask:"override"`
	ConvertEnc string   `short:"c" long:"convert" description:"convert encoding with a specified argument. e.g. c=utf8\n [utf8, utf-16, shift-jis, ...] *See Encoding spec on WHATWG" default:"false" default-mask:"original"`
}

const (
	// UNKNOWN represents file encoding cannot be detected
	UNKNOWN = "Unknown"
	// OVERRIDE represents default OutputDIr parameter
	OVERRIDE = "$OVERRIDE"
)

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

		var outPath string
		if options.OutputDir == OVERRIDE {
			outPath = file
		} else {
			outPath = options.OutputDir + "/" + filepath.Base(file)
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
			str = rep.ReplaceAllString(str, newEol)

			if options.ConvertEnc == "false" {
				// convert UTF-8 to original encoding
				converted := charsetutil.MustEncode(str, encoding)
				ioutil.WriteFile(outPath, converted, os.ModePerm)
			}
			fmt.Println(file + "\n" + "replaced eol to " + options.Replace + "\n")
		}

		if options.ConvertEnc != "false" {
			converted, err := charsetutil.Encode(str, options.ConvertEnc)
			if err != nil {
				log.Fatalln(err)
			}
			// add bom if not utf-16 before converting
			if !strings.Contains(encoding, "UTF-16") {
				converted = addBom(options, converted)
			}

			ioutil.WriteFile(outPath, converted, os.ModePerm)
			fmt.Println(file + "\n" + "converted encoding to " + options.ConvertEnc + "\n")
		}
	}
}

func addBom(options opts, data []byte) []byte {
	var converted []byte

	if strings.EqualFold(options.ConvertEnc, "utf-16") || strings.EqualFold(options.ConvertEnc, "utf-16le") {
		bom := []byte{0xFF, 0xFE}
		converted = append(data[:0], append(bom, data[0:]...)...)
	} else if strings.EqualFold(options.ConvertEnc, "utf-16be") {
		bom := []byte{0xFE, 0xFF}
		converted = append(data[:0], append(bom, data[0:]...)...)
	}

	return converted
}
