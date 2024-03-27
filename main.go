package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		if err := run(os.Stdin, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	var pathErrs []error
	for _, path := range os.Args[1:] {
		pathErrs = append(pathErrs, func() error {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			if err := run(f, os.Stdout); err != nil {
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

func run(in io.Reader, out io.Writer) error {
	var level = leveler()
	var prefix []string
	var prevLine string

	sc := bufio.NewScanner(in)
	for sc.Scan() {
		line := sc.Text()
		if strings.TrimFunc(line, isSpace) == "" {
			fmt.Fprintln(out, "\t")
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
		fmt.Fprintf(out, "%s\t%s\n",
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
		i := strings.IndexFunc(line, func(r rune) bool { return !isSpace(r) })
		if i <= 0 {
			return i
		}
		if shift == "" {
			shift = strings.Repeat(string([]rune(line)[0]), i)
		}
		level := countPrefix(line, shift)
		return level
	}
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
