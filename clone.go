package mu8

import "errors"

func Clone[G Genome](dst, src G) error {
	if dst.Len() != src.Len() {
		return errors.New("destination and source mismatch")
	}
	for i := 0; i < dst.Len(); i++ {
		dst.GetGene(i).CloneFrom(src.GetGene(i))
	}
	return nil
}
