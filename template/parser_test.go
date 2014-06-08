package template

import "testing"
import "strings"

func TestSectionRegexp(t *testing.T) {
	if !startSection.MatchString(`@sample`) {
		t.Errorf("'startSection' regexp does not work")
	}

	if !endSection.MatchString(`@sample_`) {
		t.Errorf("'endSection' regexp does not work")
	}
}

func TestCheckSyntaxSuccess(t *testing.T) {
	reader := strings.NewReader(`
@comment
Ignore this section
@comment_
`)
	if err := checkSyntax(reader); err != nil {
		t.Errorf("syntax checker mistake '")
	}
}

func TestCheckSyntaxSuccessNest(t *testing.T) {
	reader := strings.NewReader(`
@def hoge
@comment
Ignore this section
@comment_
@def_
`)
	if err := checkSyntax(reader); err != nil {
		t.Errorf("syntax checker mistake for nested case'")
	}
}

func TestCheckSyntaxFailStartSectionOnly(t *testing.T) {
	reader := strings.NewReader(`
@start
`)
	if err := checkSyntax(reader); err == nil {
		t.Errorf("syntax checker misses for only start section case")
	}
}

func TestCheckSyntaxFailEndSectionOnly(t *testing.T) {
	reader := strings.NewReader(`
@def_
`)
	if err := checkSyntax(reader); err == nil {
		t.Errorf("syntax checker miss for only end section case")
	}
}

func TestCheckSyntaxFailInvalidEndSection(t *testing.T) {
	reader := strings.NewReader(`
@foo
@bar_
`)
	if err := checkSyntax(reader); err == nil {
		t.Errorf("syntax checker miss invald end section")
	}
}
