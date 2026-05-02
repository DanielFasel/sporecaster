package main

import (
	"fmt"
	"os"

	"github.com/DanielFasel/sporecaster/internal/loader"
	"github.com/DanielFasel/sporecaster/internal/verify"
	golang "github.com/DanielFasel/sporecaster/internal/verify/golang"
	"github.com/DanielFasel/sporecaster/internal/verify/rails"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "verify":
		path := "spore.yaml"
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		s, err := loader.Load(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		checker, err := checkerFor(s.Language)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("verifying spore: %s\n\n", s.App)
		if !verify.Run(checker, s, ".") {
			os.Exit(1)
		}
	case "init":
		fmt.Fprintln(os.Stderr, "init: not yet implemented")
		os.Exit(1)
	case "inspect":
		fmt.Fprintln(os.Stderr, "inspect: not yet implemented")
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func checkerFor(language string) (verify.Checker, error) {
	switch language {
	case "golang":
		return golang.Checker{}, nil
	case "rails":
		return rails.Checker{}, nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  sporecaster verify [path]   verify codebase against spore.yaml")
	fmt.Println("  sporecaster init            scaffold a new spore-managed project")
	fmt.Println("  sporecaster inspect [path]  launch web visualizer (deferred)")
}
