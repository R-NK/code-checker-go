package main

import (
	"io/ioutil"
	"strings"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/yuin/charsetutil"
)

func TestReplaceToLF(t *testing.T) {
	args := []string{
		"-r", "LF",
		"-t", "test",
		"-o", "output",
		"-e", "cpp",
		"-e", "h",
	}
	var options opts
	args, err := flags.ParseArgs(&options, args)
	if err != nil {
		t.Error(err)
	}
	run(options)

	eol := checkEol(t, "output")
	t.Log(eol)
	if contains(eol, "CRLF") {
		t.Fatal("CRLF detected")
	}
}

func checkEol(t *testing.T, dir string) []string {
	files := listFilesByExts(dir, []string{"*"})
	if len(files) == 0 {
		t.Fatal("file not found")
	}

	var eol []string
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		encoding := detectCharEncode(content)

		b, err := charsetutil.DecodeBytes(content, encoding)
		str := string(b)

		if strings.Contains(str, "\r\n") {
			eol = append(eol, "CRLF")
		} else {
			eol = append(eol, "LF")
		}
	}

	return eol
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
