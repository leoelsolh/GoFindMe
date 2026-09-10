package main

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {

	badHash := map[string]string{
		"5381c6e3a3c89c4a7bc5f57bc3634775385f5f1d3c2d2915ddb51281708f9cbf": "InfoStealer",
	}

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		hash, err := hashFile(path)
		if err != nil {
			return err
		}

		name, found := badHash[hash]
		if found {
			fmt.Printf("Bad Hash Detected...\nPossible malicious file on disk:  %s  (%s)\n", name, path)
		}

		fmt.Printf("%s  %s\n", hash, path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum), nil
}
