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

func TestSignedMinValue(t *testing.T) {
	signedShortMin := signedMinValue(`short`, 16, 2)
	if signedShortMin != "-32768" {
		t.Errorf(`Expected: "-32768" but got %s`, signedShortMin)
	}

	signedShortMin2 := signedMinValue(`short`, 16, 1)
	if signedShortMin2 != "-32767" {
		t.Errorf(`Expected: "-32767" but got %s`, signedShortMin2)
	}

	signedLongLongMin := signedMinValue(`long long`, 8, 2)
	if signedLongLongMin != "-128LL" {
		t.Errorf(`Expected: "128" but got %s`, signedLongLongMin)
	}
}
