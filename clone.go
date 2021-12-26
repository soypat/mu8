package mu8

import "errors"

// Clone clones the Genes of src to dst. It does not
// modify src. dst should be initialized beforehand.
func Clone[G Genome](dst, src G) error {
	if dst.Len() != src.Len() {
		return errors.New("destination and source mismatch")
	}
	for i := 0; i < dst.Len(); i++ {
		dst.GetGene(i).CloneFrom(src.GetGene(i))
	}
	return nil
}
