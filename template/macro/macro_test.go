package macro

import (
	"testing"
)

func TestIgnoreExpression(t *testing.T) {
	if ok := isIgnoredExpression(`$Id`); !ok {
		t.Errorf("'$Id' should be ignorede")
	}
}

func TestWrongArgumentLength(t *testing.T) {
	m := &Macro{Name: "foo", Body: "bar", DummyArgs: []string{"a", "b"}}
	if _, err := m.Evaluate([]string{"c"}, nil); err == nil {
		t.Error("wrong argument length")
	}
}

func TestEvaluate(t *testing.T) {
	m := &Macro{Name: "foo", Body: `Hello $name`, DummyArgs: []string{"$name"}}
	val, err := m.Evaluate([]string{"John"}, nil)
	if err != nil {
		t.Error(err)
	}

	if val != "Hello John" {
		t.Error("failed macro expantion")
	}

	m = &Macro{Name: "foo", Body: "I'm $name\nI'm from $country",
		DummyArgs: []string{"$name", "$country"}}
	val, err = m.Evaluate([]string{"Tom", "Canada"}, nil)
	if err != nil {
		t.Error(err)
	}

	if val != "I'm Tom\nI'm from Canada" {
		t.Error("failed multiline macro expantion")
	}
}

func TestEvaluateWithEnv(t *testing.T) {
	m := &Macro{Name: "foo", Body: `Hello $name`}
	env := make(map[string]*Macro)
	env["$name"] = &Macro{Name: "bar", Body: "John"}

	val, err := m.Evaluate(nil, env)
	if err != nil {
		t.Error(err)
	}

	if val != "Hello John" {
		t.Errorf("failed macro expantion with env(got=%s)", val)
	}
}

func TestEvaluateWithFunctionEnv(t *testing.T) {
	m := &Macro{Name: "foo", Body: `Hello $print("John", "Smith")`}
	env := make(map[string]*Macro)
	env["$print"] = &Macro{Name: "bar", Body: "printf($family, $last)",
		DummyArgs: []string{"$family", "$last"}}

	val, err := m.Evaluate(nil, env)
	if err != nil {
		t.Error(err)
	}

	if val != `Hello printf("John", "Smith")` {
		t.Errorf("failed macro expantion with args(got=%s)", val)
	}
}
