package genes

import (
	"errors"

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
