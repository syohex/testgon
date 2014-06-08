package macro

import (
	"testing"
)

func TestConstructor(t *testing.T) {
	if _, err := New("", "", nil); err == nil {
		t.Errorf("'name' argument should not be '0'")
	}
}

func TestIgnoreExpression(t *testing.T) {
	if ok := isIgnoredExpression(`$Id`); !ok {
		t.Errorf("'$Id' should be ignorede")
	}
}

func TestWrongArgumentLength(t *testing.T) {
	m, _ := New("foo", "bar", []string{"a", "b"})
	if _, err := m.Evaluate([]string{"c"}, nil); err == nil {
		t.Error("wrong argument length")
	}
}

func TestEvaluate(t *testing.T) {
	m, _ := New("foo", `Hello $name`, []string{"$name"})
	val, err := m.Evaluate([]string{"John"}, nil)
	if err != nil {
		t.Error(err)
	}

	if val != "Hello John" {
		t.Error("failed macro expantion")
	}

	m, _ = New("foo", "I'm $name\nI'm from $country", []string{"$name", "$country"})
	val, err = m.Evaluate([]string{"Tom", "Canada"}, nil)
	if err != nil {
		t.Error(err)
	}

	if val != "I'm Tom\nI'm from Canada" {
		t.Error("failed multiline macro expantion")
	}
}

func TestEvaluateWithEnv(t *testing.T) {
	m, _ := New("foo", `Hello $name`, nil)
	env := make(map[string]*Macro)
	env["$name"], _ = New("bar", "John", nil)

	val, err := m.Evaluate(nil, env)
	if err != nil {
		t.Error(err)
	}

	if val != "Hello John" {
		t.Errorf("failed macro expantion with env(got=%s)", val)
	}
}

func TestEvaluateWithFunctionEnv(t *testing.T) {
	m, _ := New("foo", `Hello $print("John", "Smith")`, nil)
	env := make(map[string]*Macro)
	env["$print"], _ = New("bar", "printf($family, $last)", []string{"$family", "$last"})

	val, err := m.Evaluate(nil, env)
	if err != nil {
		t.Error(err)
	}

	if val != `Hello printf("John", "Smith")` {
		t.Errorf("failed macro expantion with args(got=%s)", val)
	}
}
