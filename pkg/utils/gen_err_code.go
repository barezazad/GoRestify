package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var mapErrorCodeCheck = make(map[string]bool)
var mapErrorCode = make(map[string]bool)
var arrayOfCode []string

// GenerateErrCode to generate error codes
func GenerateErrCode() {
	files, _ := getPathAllGoFiles("./")

	isChangesExist := false
	for _, v := range files {

		// open original file
		f, err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(v, "main.go") || strings.Contains(v, "gen_err_code.go") {
			continue
		}

		isExist, err := IsCodeDuplicate(f)
		if isExist {
			isChangesExist = true
			break
		}
		f.Close()
	}

	if !isChangesExist {
		return
	}

	for _, v := range files {

		// open original file
		f, err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(v, "main.go") || strings.Contains(v, "gen_err_code.go") {
			continue
		}

		// create temp file
		tmp, err := os.CreateTemp("", "replace-*")
		if err != nil {
			log.Fatal(err)
		}

		// replace while copying from f to tmp
		if err = replace(f, tmp); err != nil {
			log.Fatal(err)
		}

		// make sure the tmp file was successfully written to
		if err := tmp.Close(); err != nil {
			log.Fatal(err)
		}

		// close the file we're reading from
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}

		// overwrite the original file with the temp file
		if err := os.Rename(tmp.Name(), v); err != nil {
			log.Fatal(err)
		}

		f.Close()
		tmp.Close()
	}
}

// IsCodeDuplicate error code in file
func IsCodeDuplicate(r io.Reader) (ok bool, err error) {
	// use scanner to read line by line

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()

		getErrCode := RegexFindTerm(line, `"E[0-9]{7}"`)

		isExist := mapErrorCodeCheck[getErrCode]
		mapErrorCodeCheck[getErrCode] = true

		if getErrCode != "" && isExist {
			ok = true
			return ok, sc.Err()
		}
	}
	return ok, sc.Err()
}

// getPathAllGoFiles to find files
func getPathAllGoFiles(root string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match("*.go", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// replace error code in file
func replace(r io.Reader, w io.Writer) (err error) {
	// set nil of array
	arrayOfCode = nil
	// use scanner to read line by line
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()

		getErrCode := RegexFindTerm(line, `"E[0-9]{7}"`)
		isExist := mapErrorCode[getErrCode]
		mapErrorCode[getErrCode] = true

		if getErrCode != "" && isExist {
			newCode := fmt.Sprintf("\"%v\"", getCode())
			line = strings.ReplaceAll(line, getErrCode, newCode)
		}

		if _, err := io.WriteString(w, line+"\n"); err != nil {
			return err
		}
	}
	return sc.Err()
}

func getCode() (code string) {
	errCodePath := "err_codes.txt"

	f, err := os.Open(errCodePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		if code == "" {
			code = scanner.Text()
		}

		if code != scanner.Text() && code != "" {
			_, err := buf.Write(scanner.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(errCodePath, buf.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}

	if code == "" {
		rollbackRemoveCodes()
		log.Fatal("codes are not enough")
	}

	arrayOfCode = append(arrayOfCode, code)

	return code
}

// rollbackRemoveCodes when we have error in not nought code
func rollbackRemoveCodes() {

	f, err := os.OpenFile("err_codes.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for _, v := range arrayOfCode {
		if _, err := f.WriteString(fmt.Sprintf("%v\n", v)); err != nil {
			log.Println(err)
		}
	}
}
