package mu8

import "errors"

// Clone clones the Genes of src to dst. It does not
// modify src. dst should be initialized beforehand.
func Clone(dst, src Genome) error {
	if dst == nil {
		return errors.New("got nil destination for Clone")
	} else if src == nil {
		return errors.New("got nil source to Clone")
	} else if dst.Len() != src.Len() {
		return errors.New("destination and source mismatch")
	}
	for i := 0; i < dst.Len(); i++ {
		dst.GetGene(i).CloneFrom(src.GetGene(i))
	}
	return nil
}
