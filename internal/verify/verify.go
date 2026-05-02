package verify

import (
	"fmt"

	"github.com/DanielFasel/sporecaster/internal/loader"
	golang "github.com/DanielFasel/sporecaster/internal/verify/golang"
)

// Run verifies the spore against the codebase at root and prints a report.
// Returns true if all checks pass.
func Run(spore *loader.Spore, root string) bool {
	switch spore.Language {
	case "golang":
		return runReport(golang.Run(spore, root))
	default:
		fmt.Printf("unsupported language: %s\n", spore.Language)
		return false
	}
}

func runReport(results []golang.Result) bool {
	pass := true
	for _, r := range results {
		if r.OK {
			fmt.Printf("  ok   %s\n", r.Label)
		} else {
			pass = false
			fmt.Printf("  FAIL %s\n", r.Label)
			for _, issue := range r.Issues {
				fmt.Printf("         %s\n", issue)
			}
		}
	}
	fmt.Println()
	if pass {
		fmt.Println("spore verified: all packages and files present")
	} else {
		fmt.Println("spore verification failed")
	}
	return pass
}
