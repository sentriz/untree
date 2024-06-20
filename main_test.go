package main

import (
	"bytes"
	"path/filepath"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestFormats(t *testing.T) {
	paths, err := filepath.Glob("testdata/*.txtar")
	nilErr(t, err)

	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			runPath(t, path)
		})
	}
}

func runPath(t testing.TB, path string) {
	ar, err := txtar.ParseFile(path)
	nilErr(t, err)

	var in, expOut []byte
	for _, f := range ar.Files {
		switch f.Name {
		case "in":
			in = f.Data
		case "out":
			expOut = f.Data
		}
	}
	if len(in) == 0 || len(expOut) == 0 {
		t.Fatal("invalid test file")
	}

	var out bytes.Buffer
	err = run("", bytes.NewReader(in), &out)
	nilErr(t, err)

	if !bytes.Equal(expOut, out.Bytes()) {
		t.Errorf("\nwant:\n%s\ngot:\n%s", string(expOut), out.String())
	}
}

func nilErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got: %v", err)
	}
}
