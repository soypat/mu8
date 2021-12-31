package genes

import (
	"errors"
	"fmt"

	"github.com/soypat/mu8"
)

var (
	ErrMismatchedGeneType = errors.New("mu8.Gene argument in Splice or CloneFrom not same type as receiver")
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

var (
	_ gene[float64] = (*ConstrainedFloat)(nil)
	_ gene[float64] = (*NormalDistribution)(nil)
)
