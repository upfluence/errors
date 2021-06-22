package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

	pkg, err := exec.Command("go", "list", "-m").Output()

	if err != nil {
		l.Fatal(err.Error())
	}

	imports.LocalPrefix = string(pkg)

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
		l.Fatal(err.Error())
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
		if rewriteImport(s) {
			changed = true
		}
	}

	for _, d := range f.Decls {
		if rewriteMultiError(d) {
			changed = true
		}

		if cleanErrorFormat(d) {
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

	res, err := imports.Process(
		fname,
		buf.Bytes(),
		&imports.Options{Comments: true},
	)

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

func rewriteImport(s *ast.ImportSpec) bool {
	switch s.Path.Value {
	case "\"errors\"", "\"github.com/pkg/errors\"":
		s.Path.Value = "\"github.com/upfluence/errors\""
		return true
	}

	return false
}

var prefixErrorValueRegexp = regexp.MustCompile(`(\w+\/)*\w+\: `)

func cleanErrorFormat(d ast.Decl) bool {
	gd, ok := d.(*ast.GenDecl)

	if !ok || gd.Tok != token.VAR {
		return false
	}

	changed := false

	for _, spec := range gd.Specs {
		vs, ok := spec.(*ast.ValueSpec)

		if !ok || len(vs.Values) != 1 {
			continue
		}

		ce, ok := vs.Values[0].(*ast.CallExpr)

		if !ok || len(ce.Args) != 1 {
			continue
		}

		bl, oka := ce.Args[0].(*ast.BasicLit)
		se, oks := ce.Fun.(*ast.SelectorExpr)

		if !oka || !oks || bl.Kind != token.STRING {
			continue
		}

		x, ok := se.X.(*ast.Ident)

		if !ok || se.Sel.Name != "New" || x.Name != "errors" {
			continue
		}

		bl.Value = prefixErrorValueRegexp.ReplaceAllString(bl.Value, "")

		changed = true
	}

	return changed
}

func rewriteMultiError(d ast.Decl) bool {
	fd, ok := d.(*ast.FuncDecl)

	if !ok {
		return false
	}

	var changed bool

	for _, stmt := range fd.Body.List {
		if rewriteStmtMultiError(stmt) {
			changed = true
		}
	}

	return changed
}

func rewriteStmtMultiError(stmt ast.Stmt) bool {
	var changed bool

	switch sstmt := stmt.(type) {
	case *ast.DeclStmt:
		return rewriteMultiError(sstmt.Decl)
	case *ast.AssignStmt:
		for _, r := range sstmt.Rhs {
			if rewriteExprMultiError(r) {
				changed = true
			}
		}
	case *ast.IfStmt:
		return rewriteStmtMultiError(sstmt.Init)
	case *ast.ReturnStmt:
		for _, r := range sstmt.Results {
			if rewriteExprMultiError(r) {
				changed = true
			}
		}
	}

	return changed
}

func rewriteExprMultiError(expr ast.Expr) bool {
	switch texpr := expr.(type) {
	case *ast.SelectorExpr:
		pkg, ok := texpr.X.(*ast.Ident)

		if !ok || pkg.Name != "multierror" {
			return false
		}

		pkg.Name = "errors"

		if texpr.Sel.Name == "Wrap" {
			texpr.Sel.Name = "WrapErrors"
		}

		return true
	case *ast.CallExpr:
		return rewriteExprMultiError(texpr.Fun)
	}

	return false
}
