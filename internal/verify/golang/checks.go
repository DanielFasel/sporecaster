package golang

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/DanielFasel/sporecaster/internal/spore"
	"github.com/DanielFasel/sporecaster/internal/verify"
)

type Checker struct{}

func (c Checker) Run(s *spore.Spore, root string) []verify.Result {
	var results []verify.Result

	// ── Zoom 1: skeleton ────────────────────────────────────────────────

	corePath := filepath.Join(root, "cmd", s.Core.Name)
	results = append(results, checkDir(s.Core.Name+" (core)", corePath, s.Core.Files))
	if dirExists(corePath) {
		results = append(results, checkPkgDecl(s.Core.Name+" (core)", corePath, "main"))
	}

	for _, pkg := range s.Packages {
		pkgPath := filepath.Join(root, "internal", filepath.FromSlash(pkg.Name))
		results = append(results, checkDir(pkg.Name, pkgPath, pkg.Files))
		if dirExists(pkgPath) {
			results = append(results, checkPkgDecl(pkg.Name, pkgPath, leafName(pkg.Name)))
		}
	}

	// ── Zoom 2: connections ─────────────────────────────────────────────

	modulePath, err := readModulePath(root)
	if err != nil {
		results = append(results, verify.Result{
			Zoom:   2,
			Label:  "module path",
			OK:     false,
			Issues: []string{err.Error()},
		})
		return results
	}

	// import graph
	if dirExists(corePath) {
		results = append(results, checkImports(s.Core.Name+" (core)", corePath, modulePath, s.Core.Imports))
	}
	for _, pkg := range s.Packages {
		pkgPath := filepath.Join(root, "internal", filepath.FromSlash(pkg.Name))
		if dirExists(pkgPath) {
			results = append(results, checkImports(pkg.Name, pkgPath, modulePath, pkg.Imports))
		}
	}

	// CLI channel commands handled in core
	if dirExists(corePath) {
		results = append(results, checkCLICommands(s.Core.Name+" (core)", corePath, s.Channels)...)
	}

	// error handling: os.Exit only in main
	if s.ErrorHandling.TerminatesAt == "main" {
		for _, pkg := range s.Packages {
			pkgPath := filepath.Join(root, "internal", filepath.FromSlash(pkg.Name))
			if dirExists(pkgPath) {
				results = append(results, checkNoOsExit(pkg.Name, pkgPath))
			}
		}
	}

	// ── Zoom 3: exports ────────────────────────────────────────────────

	for _, pkg := range s.Packages {
		pkgPath := filepath.Join(root, "internal", filepath.FromSlash(pkg.Name))
		if dirExists(pkgPath) && len(pkg.Exports) > 0 {
			results = append(results, checkExports(pkg.Name, pkgPath, pkg.Exports))
		}
	}

	// error handling: sentinel errors must be exported
	if s.ErrorHandling.SentinelsAreExported {
		if dirExists(corePath) {
			results = append(results, checkExportedSentinels(s.Core.Name+" (core)", corePath))
		}
		for _, pkg := range s.Packages {
			pkgPath := filepath.Join(root, "internal", filepath.FromSlash(pkg.Name))
			if dirExists(pkgPath) {
				results = append(results, checkExportedSentinels(pkg.Name, pkgPath))
			}
		}
	}

	return results
}

// ── zoom 1 checks ────────────────────────────────────────────────────────────

func checkDir(label, dir string, files []spore.File) verify.Result {
	r := verify.Result{Zoom: 1, Label: label, OK: true}

	if !dirExists(dir) {
		r.OK = false
		r.Issues = append(r.Issues, fmt.Sprintf("directory not found: %s", dir))
		return r
	}

	for _, f := range files {
		if _, err := os.Stat(filepath.Join(dir, f.Name)); os.IsNotExist(err) {
			r.OK = false
			r.Issues = append(r.Issues, fmt.Sprintf("missing file: %s", f.Name))
		}
	}

	return r
}

func checkPkgDecl(label, dir, want string) verify.Result {
	r := verify.Result{Zoom: 1, Label: label + ": package declaration", OK: true}

	entries, err := os.ReadDir(dir)
	if err != nil {
		r.OK = false
		r.Issues = append(r.Issues, fmt.Sprintf("cannot read directory: %v", err))
		return r
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		f, err := os.Open(filepath.Join(dir, e.Name()))
		if err != nil {
			r.OK = false
			r.Issues = append(r.Issues, fmt.Sprintf("%s: cannot open: %v", e.Name(), err))
			continue
		}
		declared, found := readPkgDecl(f)
		f.Close()

		if !found {
			r.OK = false
			r.Issues = append(r.Issues, fmt.Sprintf("%s: no package declaration", e.Name()))
			continue
		}
		if declared != want {
			r.OK = false
			r.Issues = append(r.Issues, fmt.Sprintf("%s: package %s, want %s", e.Name(), declared, want))
		}
	}

	return r
}

// ── zoom 2 checks ────────────────────────────────────────────────────────────

func checkImports(label, dir, modulePath string, declared []string) verify.Result {
	r := verify.Result{Zoom: 2, Label: label + ": imports", OK: true}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ImportsOnly)
	if err != nil {
		r.OK = false
		r.Issues = append(r.Issues, fmt.Sprintf("parse error: %v", err))
		return r
	}

	prefix := modulePath + "/internal/"
	actual := map[string]bool{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, imp := range file.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				if strings.HasPrefix(path, prefix) {
					actual[strings.TrimPrefix(path, prefix)] = true
				}
			}
		}
	}

	want := map[string]bool{}
	for _, imp := range declared {
		want[imp] = true
	}

	for imp := range want {
		if !actual[imp] {
			r.OK = false
			r.Issues = append(r.Issues, fmt.Sprintf("declared but missing in code: %s", imp))
		}
	}
	for imp := range actual {
		if !want[imp] {
			r.OK = false
			r.Issues = append(r.Issues, fmt.Sprintf("present in code but undeclared: %s", imp))
		}
	}

	return r
}

func checkCLICommands(label, dir string, channels []spore.Channel) []verify.Result {
	var results []verify.Result
	for _, ch := range channels {
		if ch.Type != "cli" {
			continue
		}
		for _, cmd := range ch.Commands {
			r := verify.Result{Zoom: 2, Label: fmt.Sprintf("%s: cli command %q", label, cmd.Name), OK: true}
			if !caseHandled(dir, cmd.Name) {
				r.OK = false
				r.Issues = append(r.Issues, fmt.Sprintf("no case %q found in switch", cmd.Name))
			}
			results = append(results, r)
		}
	}
	return results
}

func checkNoOsExit(label, dir string) verify.Result {
	r := verify.Result{Zoom: 2, Label: label + ": no os.Exit", OK: true}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		r.OK = false
		r.Issues = append(r.Issues, fmt.Sprintf("parse error: %v", err))
		return r
	}

	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				id, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}
				if id.Name == "os" && sel.Sel.Name == "Exit" {
					r.OK = false
					r.Issues = append(r.Issues, fmt.Sprintf("%s: os.Exit called", filepath.Base(filename)))
				}
				return true
			})
		}
	}

	return r
}

func checkExportedSentinels(label, dir string) verify.Result {
	r := verify.Result{Zoom: 2, Label: label + ": exported sentinels", OK: true}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		r.OK = false
		r.Issues = append(r.Issues, fmt.Sprintf("parse error: %v", err))
		return r
	}

	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			for _, decl := range file.Decls {
				gd, ok := decl.(*ast.GenDecl)
				if !ok || gd.Tok != token.VAR {
					continue
				}
				for _, spec := range gd.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for i, val := range vs.Values {
						if !isSentinelCall(val) {
							continue
						}
						if i >= len(vs.Names) {
							continue
						}
						name := vs.Names[i].Name
						if !ast.IsExported(name) {
							r.OK = false
							r.Issues = append(r.Issues, fmt.Sprintf("%s: sentinel %s is unexported", filepath.Base(filename), name))
						}
					}
				}
			}
		}
	}

	return r
}

// ── zoom 3 checks ────────────────────────────────────────────────────────────

func checkExports(label, dir string, declared []spore.Export) verify.Result {
	r := verify.Result{Zoom: 3, Label: label + ": exports", OK: true}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		r.OK = false
		r.Issues = append(r.Issues, fmt.Sprintf("parse error: %v", err))
		return r
	}

	// collect exported top-level symbols from code
	actual := map[string]string{} // name → kind
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Recv == nil && d.Name.IsExported() {
						actual[d.Name.Name] = "func"
					}
				case *ast.GenDecl:
					if d.Tok != token.TYPE {
						continue
					}
					for _, spec := range d.Specs {
						ts, ok := spec.(*ast.TypeSpec)
						if !ok || !ts.Name.IsExported() {
							continue
						}
						switch ts.Type.(type) {
						case *ast.StructType:
							actual[ts.Name.Name] = "struct"
						case *ast.InterfaceType:
							actual[ts.Name.Name] = "interface"
						default:
							actual[ts.Name.Name] = "type"
						}
					}
				}
			}
		}
	}

	want := map[string]string{}
	for _, e := range declared {
		want[e.Name] = e.Kind
	}

	for name, wantKind := range want {
		gotKind, exists := actual[name]
		if !exists {
			r.OK = false
			r.Issues = append(r.Issues, fmt.Sprintf("%s: not found (want %s)", name, wantKind))
		} else if gotKind != wantKind {
			r.OK = false
			r.Issues = append(r.Issues, fmt.Sprintf("%s: is %s, want %s", name, gotKind, wantKind))
		}
	}

	for name := range actual {
		if _, declared := want[name]; !declared {
			r.OK = false
			r.Issues = append(r.Issues, fmt.Sprintf("%s: exported but not declared in spec", name))
		}
	}

	return r
}

// ── helpers ──────────────────────────────────────────────────────────────────

func caseHandled(dir, cmd string) bool {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return false
	}
	found := false
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if found {
					return false
				}
				cc, ok := n.(*ast.CaseClause)
				if !ok {
					return true
				}
				for _, expr := range cc.List {
					lit, ok := expr.(*ast.BasicLit)
					if ok && lit.Kind == token.STRING && strings.Trim(lit.Value, `"`) == cmd {
						found = true
					}
				}
				return true
			})
		}
	}
	return found
}

func isSentinelCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return (id.Name == "errors" && sel.Sel.Name == "New") ||
		(id.Name == "fmt" && sel.Sel.Name == "Errorf")
}

func readModulePath(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("reading go.mod: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "module ") {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "module ")), nil
		}
	}
	return "", fmt.Errorf("no module directive in go.mod")
}

func readPkgDecl(f *os.File) (string, bool) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "package ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1], true
			}
		}
	}
	return "", false
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func leafName(name string) string {
	parts := strings.Split(name, "/")
	return parts[len(parts)-1]
}
