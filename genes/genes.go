package genes

import (
	"constraints"
	"errors"
	"fmt"

	"github.com/soypat/mu8"
)

var (
	ErrMismatchedGeneType = errors.New("mu8.Gene argument in Splice or CloneFrom not same type as receiver")
	errStartOutOfBounds   = errors.New("start value should be contained within bounds [min,max] for Contrained types")
	errBadConstraints     = errors.New("min should be less than max for Constrained types and not equal for int gene types")
)

// Helper function for casting interfaces.
func castGene[T mu8.Gene](gene mu8.Gene) T {
	g, ok := gene.(T)
	if !ok {
		panic(ErrMismatchedGeneType.Error())
	}
	return g
}

// Compile time check of internal interface implementation

type gene[T any] interface {
	mu8.Gene
	fmt.Stringer
	Value() T
	SetValue(v T)
}

// Compile-time checks of interface implementation.
var (
	_ gene[float64] = (*ConstrainedFloat)(nil)
	_ gene[float64] = (*NormalDistribution)(nil)
	_ gene[float64] = (*ConstrainedNormalDistr)(nil)
	_ gene[int]     = (*ConstrainedInt)(nil)
)

func max[T constraints.Integer](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T constraints.Integer](a, b T) T {
	if a < b {
		return a
	}
	return b
}
