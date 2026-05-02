package rails

import (
	"github.com/DanielFasel/sporecaster/internal/spore"
	"github.com/DanielFasel/sporecaster/internal/verify"
)

type Checker struct{}

func (c Checker) Run(s *spore.Spore, root string) []verify.Result {
	return nil
}
