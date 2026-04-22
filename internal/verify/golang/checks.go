package golang

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/you/sporecast/internal/loader"
)

type Result struct {
	Label  string
	OK     bool
	Issues []string
}

// Run checks that every package and file declared in the spore exists on disk.
func Run(spore *loader.Spore, root string) []Result {
	var results []Result

	corePath := filepath.Join(root, "cmd", spore.Core.Name)
	results = append(results, checkPackage(spore.Core.Name+" (core)", corePath, spore.Core.Files))

	for _, pkg := range spore.Packages {
		pkgPath := filepath.Join(root, "internal", pkg.Name)
		results = append(results, checkPackage(pkg.Name, pkgPath, pkg.Files))
	}

	for _, sub := range spore.SubPackages {
		subPath := filepath.Join(root, "internal", filepath.FromSlash(sub.Key))
		results = append(results, checkPackage(sub.Key, subPath, sub.Files))
	}

	return results
}

func checkPackage(label, dir string, files []loader.File) Result {
	r := Result{Label: label, OK: true}

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
