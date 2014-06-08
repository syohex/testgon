package template

import (
	"testing"
)

func TestFileSectionRegexp(t *testing.T) {
	input := `@file >>fn???_extern.c $macro1() @file_`
	if !fileSectionRegexp.MatchString(input) {
		t.Errorf("Can't match to '%s'", input)
	}
}
