package golang

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DanielFasel/sporecaster/internal/spore"
	"github.com/DanielFasel/sporecaster/internal/verify"
)

type Checker struct{}

func (c Checker) Run(s *spore.Spore, root string) []verify.Result {
	var results []verify.Result

	corePath := filepath.Join(root, "cmd", s.Core.Name)
	results = append(results, checkPackage(s.Core.Name+" (core)", corePath, s.Core.Files))

	for _, pkg := range s.Packages {
		pkgPath := filepath.Join(root, "internal", filepath.FromSlash(pkg.Name))
		results = append(results, checkPackage(pkg.Name, pkgPath, pkg.Files))
	}

	return results
}

func checkPackage(label, dir string, files []spore.File) verify.Result {
	r := verify.Result{Label: label, OK: true}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
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
