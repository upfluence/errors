package main

import (
	"bytes"
	"flag"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
)

func main() {
	var (
		path   string
		dryRun bool

		l = log.New(os.Stderr, "", log.LstdFlags)
	)

	flag.StringVar(&path, "path", ".", "directory to migrate")
	flag.BoolVar(&dryRun, "dry-run", false, "dry run")

	flag.Parse()

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		switch path {
		case ".git", "vendor":
			return filepath.SkipDir
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		return parseFile(path, l, dryRun)
	})

	if err != nil {
		l.Println(err.Error())
	}
}

func parseFile(fname string, l *log.Logger, dryRun bool) error {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fname, nil, 0)

	if err != nil {
		return err
	}

	var changed bool

	for _, s := range f.Imports {
		switch s.Path.Value {
		case "\"errors\"", "\"github.com/pkg/errors\"":
			s.Path.Value = "\"github.com/upfluence/errors\""
			changed = true
		}
	}

	if !changed {
		return nil
	}

	l.Println(fname)

	var (
		w io.Writer
		c io.Closer
	)

	if dryRun {
		w = l.Writer()
	} else {
		file, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0)

		if err != nil {
			return err
		}

		w = file
		c = file
	}

	var buf bytes.Buffer

	if err := format.Node(&buf, fset, f); err != nil {
		return err
	}

	res, err := imports.Process(fname, buf.Bytes(), &imports.Options{FormatOnly: true})

	if err != nil {
		return err
	}

	if _, err := w.Write(res); err != nil {
		return err
	}

	if c != nil {
		if err := c.Close(); err != nil {
			return err
		}
	}

	return nil
}
