package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func main() {
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
