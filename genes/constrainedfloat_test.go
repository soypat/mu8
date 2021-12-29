package genes

import (
	"fmt"
	"testing"
)

func TestConstrainedFloatFormat(t *testing.T) {
	cf := NewConstrainedFloat(.5, 0, 1)
	str := fmt.Sprintf("%s %.2f %v", cf, cf, cf)
	t.Error(str)
}
