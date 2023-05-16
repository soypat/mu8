package genes

import (
	"fmt"
	"testing"
)

func TestConstrainedFloatFormat(t *testing.T) {
	cf := NewConstrainedFloat(.5, 0, 1)
	str1 := fmt.Sprintf("%s", cf)
	str2 := fmt.Sprintf("%f", cf)
	if str1 != str2 || str1 != fmt.Sprintf("%f", 0.5) {
		t.Error("bad format")
	}
}
