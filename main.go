package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func listFiles(pth *string, filetype *string) ([]string, error) {
	var csFiles []string
	err := filepath.Walk(*pth, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() &&
			filepath.Ext(path) == *filetype {
			csFiles = append(csFiles, path)
		}

		return nil
	})

	if err != nil {
		return csFiles, err
	}

	return csFiles, nil
}

func searchPhrase(pth *string, phrase *string) (bool, int, error) {
	f, err := os.Open(*pth)
	if err != nil {
		return false, 0, err
	}
	defer f.Close()

	toSearch := strings.ToLower(*phrase)

	scn := bufio.NewScanner(f)
	cnt := 0
	for scn.Scan() {
		cnt += 1
		if !strings.Contains(strings.ToLower(scn.Text()), toSearch) {
			continue
		}

		return true, cnt, nil
	}

	if err = scn.Err(); err != nil {
		return false, cnt, err
	}

	return false, 0, err
}

func searchPhraseAll(pth *string, filetype *string, phrase *string) ([]string, error) {
	var matchFiles []string

	csFiles, err := listFiles(pth, filetype)
	if err != nil {
		return matchFiles, err
	}

	for _, fName := range csFiles {
		hasPhrase, lineNum, err := searchPhrase(&fName, phrase)
		if err != nil {
			return matchFiles, err
		}

		if !hasPhrase {
			continue
		}
		var sBuff bytes.Buffer
		sBuff.WriteString(fName)
		sBuff.WriteString(" at line: ")
		sBuff.WriteString(strconv.Itoa(lineNum))
		matchFiles = append(matchFiles, sBuff.String())
	}

	return matchFiles, nil
}

func main() {
	phrase := flag.String("phrase", "shoot", "# phrase to search recursivley")
	filetype := flag.String("filetype", ".cs", "# file types to serach")
	filepath := flag.String("filepath", "./data/", "# dir to search in")
	flag.Parse()
	matchFiles, err := searchPhraseAll(filepath, filetype, phrase)
	if err != nil {
		log.Fatal(err)
	}

	for _, fName := range matchFiles {
		fmt.Println(fName)
	}
}
