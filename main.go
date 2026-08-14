package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	withFilepaths := flag.Bool("paths", false, "include filepaths in output")
	withVim := flag.Bool("vim", false, "output as path:line:col: for quickfix / grepprg")
	flag.Parse()

	prefixer := func(path string) prefixFunc {
		switch {
		case *withVim:
			return func(lineNum, col int) string { return fmt.Sprintf("%s:%d:%d:", path, lineNum, col) }
		case *withFilepaths:
			return func(int, int) string { return path + "\t" }
		}
		return noPrefix
	}

	if flag.NArg() == 0 {
		if err := run(prefixer(""), os.Stdin, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	var pathErrs []error
	for _, path := range flag.Args() {
		prefix := prefixer(path)
		pathErrs = append(pathErrs, func() error {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			if err := run(prefix, f, os.Stdout); err != nil {
				return fmt.Errorf("run: %w", err)
			}
			return nil
		}())
	}
	if err := errors.Join(pathErrs...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// prefixFunc renders the per-line prefix for the 1-indexed line number and the
// column of its first non-space byte.
type prefixFunc func(lineNum, col int) string

func noPrefix(int, int) string { return "" }

func run(linePrefix prefixFunc, in io.Reader, out io.Writer) error {
	var level = leveler()
	var prefix []string
	var prevLine string
	var lineNum int

	sc := bufio.NewScanner(in)
	for sc.Scan() {
		line := sc.Text()
		lineNum++
		if strings.TrimFunc(line, isSpace) == "" {
			fmt.Fprintf(out, "%s\t\n", linePrefix(lineNum, 1))
			continue
		}
		if l := level(line); l > len(prefix) {
			for len(prefix) < l-1 { // pad up to l in case we indented more than once
				prefix = append(prefix, "")
			}
			prefix = append(prefix, strings.TrimSpace(prevLine))
		} else if l < len(prefix) {
			prefix = prefix[:l]
		}
		fmt.Fprintf(out, "%s%s\t%s\n",
			linePrefix(lineNum, indentWidth(line)+1),
			strings.TrimSpace(strings.Join(prefix, " ")),
			strings.ReplaceAll(line, "\t", "    "),
		)
		prevLine = line
	}
	return sc.Err()
}

func leveler() func(line string) int {
	var shift string
	return func(line string) int {
		i := indentWidth(line)
		if i == 0 {
			return 0
		}
		if shift == "" {
			shift = strings.Repeat(string([]rune(line)[0]), i)
		}
		level := countPrefix(line, shift)
		return level
	}
}

func indentWidth(line string) int {
	if i := strings.IndexFunc(line, func(r rune) bool { return !isSpace(r) }); i > 0 {
		return i
	}
	return 0
}

func countPrefix(line, p string) int {
	if len(p) == 0 {
		return 0
	}
	var count int
	for i := 0; i+len(p)-1 < len(line) && line[i:i+len(p)] == p; i += len(p) {
		count++
	}
	return count
}

// custom isSpace for our considered indent chars
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}
