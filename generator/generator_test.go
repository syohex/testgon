package generator

import (
	"testing"
)

func TestSignedMaxValue(t *testing.T) {
	signedCharMax := signedMaxValue(`char`, 8)
	if signedCharMax != "127" {
		t.Errorf(`Expected: "127" but got %s`, signedCharMax)
	}

	signedLongMax := signedMaxValue(`long`, 32)
	if signedLongMax != "2147483647L" {
		t.Errorf(`Expected: "2147483647L" but got %s`, signedLongMax)
	}
}
