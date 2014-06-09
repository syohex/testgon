package generator

import (
	"testing"
)

func TestSignedMaxValue(t *testing.T) {
	signedCharMax := signedMaxValue(`char`, 8)
	if signedCharMax != "127" {
		t.Errorf("Expected: 127 but got %s", signedCharMax)
	}
}
