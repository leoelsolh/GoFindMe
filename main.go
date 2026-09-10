package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	badHash, err := loadSignatures("signatures.txt")
	if err != nil {
		fmt.Println("Could not load signatures:", err)
		return
	}

	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
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

func loadSignatures(path string) (map[string]string, error) {
	sigs := map[string]string{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		hash, name, found := strings.Cut(line, ",")
		if !found {
			continue
		}
		sigs[hash] = name
	}

	return sigs, nil
}
