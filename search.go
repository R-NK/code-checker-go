package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func listFilesByExts(dir string, exts []string) []string {
	paths := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
			return err
		}
		for _, ext := range exts {
			// ディレクトリの場合なにもしない
			fi, _ := os.Stat(path)
			if fi.IsDir() {
				continue
			}
			// 拡張子が.から始まらない場合先頭に付与
			if !strings.Contains(ext, ".") {
				ext = "." + ext
			}
			// ワイルドカードの場合そのまま追加
			if ext == filepath.Ext(path) || ext == ".*" {
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
