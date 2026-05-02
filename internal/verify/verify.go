package verify

import (
	"fmt"
	"sort"

	"github.com/DanielFasel/sporecaster/internal/spore"
)

type Result struct {
	Zoom   int
	Label  string
	OK     bool
	Issues []string
}

type Checker interface {
	Run(s *spore.Spore, root string) []Result
}

func Run(checker Checker, s *spore.Spore, root string) bool {
	results := checker.Run(s, root)

	// collect zoom levels in the order they appear, then sort
	seen := map[int]bool{}
	var zoomOrder []int
	for _, r := range results {
		if !seen[r.Zoom] {
			seen[r.Zoom] = true
			zoomOrder = append(zoomOrder, r.Zoom)
		}
	}
	sort.Ints(zoomOrder)

	pass := true
	for _, z := range zoomOrder {
		fmt.Printf("Zoom %d — %s\n", z, zoomLabel(z))
		for _, r := range results {
			if r.Zoom != z {
				continue
			}
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
	}

	if pass {
		fmt.Println("spore verified: all checks passed")
	} else {
		fmt.Println("spore verification failed")
	}
	return pass
}

func zoomLabel(z int) string {
	switch z {
	case 1:
		return "skeleton"
	case 2:
		return "connections"
	case 3:
		return "exports"
	default:
		return fmt.Sprintf("zoom %d", z)
	}
}
