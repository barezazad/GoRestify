package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GenerateErrCode finds duplicate E####### error codes in .go files and
// replaces later duplicates by bumping the numeric part (E1146731 → E1146732 → …)
// until the code is unused in the project. Does not use err_codes.txt.
func GenerateErrCode() {
	files, err := getPathAllGoFiles("./")
	if err != nil {
		log.Fatal(err)
	}

	used := collectErrCodes(files)
	if !hasDuplicateErrCodes(files) {
		return
	}

	// Track first-seen codes while rewriting; bump later collisions.
	seen := make(map[string]bool)

	for _, path := range files {
		if shouldSkipErrCodeFile(path) {
			continue
		}

		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		tmp, err := os.CreateTemp("", "replace-*")
		if err != nil {
			f.Close()
			log.Fatal(err)
		}

		if err = rewriteErrCodes(f, tmp, seen, used); err != nil {
			f.Close()
			tmp.Close()
			log.Fatal(err)
		}

		if err := f.Close(); err != nil {
			tmp.Close()
			log.Fatal(err)
		}
		if err := tmp.Close(); err != nil {
			log.Fatal(err)
		}
		if err := os.Rename(tmp.Name(), path); err != nil {
			log.Fatal(err)
		}
	}
}

func shouldSkipErrCodeFile(path string) bool {
	base := filepath.Base(path)
	return base == "gen_err_code.go" || strings.HasSuffix(base, "main.go") ||
		strings.Contains(path, "main.admin.go") || strings.Contains(path, "main.user.go")
}

func collectErrCodes(files []string) map[string]bool {
	used := make(map[string]bool)
	for _, path := range files {
		if shouldSkipErrCodeFile(path) {
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if code := extractErrCode(sc.Text()); code != "" {
				used[code] = true
			}
		}
		if err := sc.Err(); err != nil {
			f.Close()
			log.Fatal(err)
		}
		f.Close()
	}
	return used
}

func hasDuplicateErrCodes(files []string) bool {
	seen := make(map[string]bool)
	for _, path := range files {
		if shouldSkipErrCodeFile(path) {
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			code := extractErrCode(sc.Text())
			if code == "" {
				continue
			}
			if seen[code] {
				f.Close()
				return true
			}
			seen[code] = true
		}
		if err := sc.Err(); err != nil {
			f.Close()
			log.Fatal(err)
		}
		f.Close()
	}
	return false
}

func rewriteErrCodes(r io.Reader, w io.Writer, seen, used map[string]bool) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		quoted := RegexFindTerm(line, `"E[0-9]{7}"`)
		if quoted != "" {
			code := strings.Trim(quoted, `"`)
			if seen[code] {
				next := nextUnusedErrCode(code, used)
				line = strings.ReplaceAll(line, quoted, `"`+next+`"`)
				seen[next] = true
			} else {
				seen[code] = true
			}
		}
		if _, err := io.WriteString(w, line+"\n"); err != nil {
			return err
		}
	}
	return sc.Err()
}

func extractErrCode(line string) string {
	quoted := RegexFindTerm(line, `"E[0-9]{7}"`)
	if quoted == "" {
		return ""
	}
	return strings.Trim(quoted, `"`)
}

// nextUnusedErrCode bumps E1146731 → E1146732 → … until unused, then marks it used.
func nextUnusedErrCode(code string, used map[string]bool) string {
	if len(code) != 8 || code[0] != 'E' {
		log.Fatalf("invalid error code format: %s", code)
	}

	n, err := strconv.ParseUint(code[1:], 10, 64)
	if err != nil {
		log.Fatalf("invalid error code number: %s", code)
	}

	for {
		n++
		if n > 9999999 {
			log.Fatalf("no unused error code left after %s", code)
		}
		candidate := fmt.Sprintf("E%07d", n)
		if !used[candidate] {
			used[candidate] = true
			return candidate
		}
	}
}

func getPathAllGoFiles(root string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			name := info.Name()
			if name == "vendor" || name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if matched, err := filepath.Match("*.go", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	return matches, err
}
