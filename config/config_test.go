package config

import (
	"testing"
)

const sample1 = `
{
  "compiler":  "gcc",
  "simulator": "gdb",
  "c_flags":     [ "-g", "-Dunix" ],
  "ld_flags":    [ "-lm" ],
  "options":     [ "-O0", "-O2" ],
  "testdir":     "testsuite",
  "compile_only": true,

  "size": {
       "char": 8,
       "short": 16,
       "int": 32,
       "long": 64,
       "pointer" : 64
  },

  "timeout": 5,
  "parallels": 2,
  "color": true,

  "lang": "c++",
  "expect": "#OK#",
  "complement": 1,
  "option_separator": "--",
  "has_printf": false
}
`
func Test1(t *testing.T) {
	conf, err := parseBytes([]byte(sample1))
	if err != nil {
		t.Error(err)
	}

	if conf.Compiler != "gcc" {
		t.Error("'compiler' parameter is not set")
	}

	if conf.Simulator != "gdb" {
		t.Error("'simulator' parameter is not set")
	}

	if !(conf.CFlags[0] == "-g" && conf.CFlags[1] == "-Dunix") {
		t.Error("'c_flags' parameter is not set")
	}

	if conf.LDFlags[0] != "-lm" {
		t.Error("'ld_flags' parameter is not set")
	}

	if !(conf.Options[0] == "-O0" && conf.Options[1] == "-O2") {
		t.Error("'options' parameter is not set")
	}

	if conf.TestDir != "testsuite" {
		t.Error("'test' parameter is not set")
	}

	if conf.CompileOnly != true {
		t.Error("'compile_only' parameter is not set")
	}

	if conf.Size.Char != 8 {
		t.Error("'size.char' parameter is not set")
	}

	if conf.Size.Short != 16 {
		t.Error("'size.short' parameter is not set")
	}

	if conf.Size.Int != 32 {
		t.Error("'size.int' parameter is not set")
	}

	if conf.Size.Long != 64 {
		t.Error("'size.long' parameter is not set")
	}

	if conf.Size.Pointer != 64 {
		t.Error("'size.pointer' parameter is not set")
	}

	if conf.Timeout != 5 {
		t.Errorf("'timeout' parameter is not set(timeout=%d)", conf.Timeout)
	}

	if conf.Parallels != 2 {
		t.Error("'parallels' parameter is not set")
	}

	if conf.Color != true {
		t.Error("'color' parameter is not set")
	}

	if conf.Lang != "c++" {
		t.Error("'lang' parameter is not set")
	}

	if conf.Expect != "#OK#" {
		t.Error("'expect' parameter is not set")
	}

	if conf.Complement != 1 {
		t.Error("'complement' parameter is not set")
	}

	if conf.OptionSeparator != "--" {
		t.Error("'option_separator' parameter is not set")
	}

	if conf.HasPrintf != false {
		t.Errorf("'has_printf' parameter is not set(got=%v)", conf.HasPrintf)
	}
}

func TestDefaultValue(t *testing.T) {
	jsonStr := `
{
    "compiler": "gcc", "testdir": "testsuite",
    "size": { "char": 8, "short": 16, "int": 32, "long": 64 }
}
`
	conf, err := parseBytes([]byte(jsonStr))
	if err != nil {
		t.Error(err)
	}

	if conf.Lang != "c" {
		t.Errorf("Default 'Lang' value is '%s' not 'c'", conf.Lang)
	}

	if conf.Parallels != 1 {
		t.Errorf("Default 'Parallels' value is '%d' not '1'", conf.Parallels)
	}

	if conf.Timeout != 10 {
		t.Errorf("Default 'Timeout' value is '%d' not '10'", conf.Timeout)
	}

	if conf.Expect != "@OK@" {
		t.Errorf("Default 'Expect' value is '%d' not '@OK@'", conf.Expect)
	}

	if conf.Complement != 2 {
		t.Errorf("Default 'Complement' value is '%d' not '2'", conf.Complement)
	}

	if conf.HasPrintf != true {
		t.Errorf("Default 'HasPrintf' value is '%v' not 'true'", conf.HasPrintf)
	}
}

func TestNotFoundCompiler(t *testing.T) {
	jsonStr := `
{
    "compiler": "not_found_compiler",
    "testdir":  "testsuite",
    "size": { "char": 8, "short": 16, "int": 32, "long": 64 }
}
`
	if _, err := parseBytes([]byte(jsonStr)); err == nil {
		t.Error("compiler not found but error is not returned")
	}
}

func TestNotFoundSimulator(t *testing.T) {
	jsonStr := `
{
    "compiler": "cc", "testdir": "testsuite",
    "simulator": "not_found_simulator"
    "size": { "char": 8, "short": 16, "int": 32, "long": 64, }
}
`
	if _, err := parseBytes([]byte(jsonStr)); err == nil {
		t.Error("simulator not found but error is not returned")
	}
}

func TestMandatoryParameters(t *testing.T) {
	jsonStr := `
{
    "testdir": "testsuite",
    "size": { "char": 8, "short": 16, "int": 32, "long": 64, "pointer": 64 }
}
`
	if _, err := parseBytes([]byte(jsonStr)); err == nil {
		t.Error("'compiler' parameter not found but error is not returned")
	}

	jsonStr = `
{
  "compiler": "cc",
  "size": { "char": 8, "short": 16, "int": 32, "long": 64, "pointer": 64 }
}
`
	if _, err := parseBytes([]byte(jsonStr)); err == nil {
		t.Error("'testdir' parameter not found but error is not returned")
	}

	jsonStr = `
{
  "compiler": "cc", "testdir": "testsuite",
}
`
	if _, err := parseBytes([]byte(jsonStr)); err == nil {
		t.Error("'size' parameter not found but error is not returned")
	}
}

func TestSizeParameter(t *testing.T) {
	jsonStr := `
{
  "compiler": "cc", "testdir": "testsuite",
  "size": { "char": 8, "short": 16, "int": 32, "long": 64 }
}
`
	if _, err := parseBytes([]byte(jsonStr)); err != nil {
		t.Error("'pointer' parameter can be omitted")
	}

	jsonStr = `
{
  "compiler": "cc", "testdir": "testsuite",
  "size": { "short": 16, "int": 32, "long": 64 }
}
`
	if _, err := parseBytes([]byte(jsonStr)); err == nil {
		t.Error("'char' parameter can not be omitted")
	}

	jsonStr = `
{
  "compiler": "cc", "testdir": "testsuite",
  "size": { "char": 8, "int": 32, "long": 64 }
}
`
	if _, err := parseBytes([]byte(jsonStr)); err == nil {
		t.Error("'short' parameter can not be omitted")
	}

	jsonStr = `
{
  "compiler": "cc", "testdir": "testsuite",
  "size": { "char": 8, "short": 32, "long": 64 }
}
`
	if _, err := parseBytes([]byte(jsonStr)); err == nil {
		t.Error("'int' parameter can not be omitted")
	}

	jsonStr = `
{
  "compiler": "cc", "testdir": "testsuite",
  "size": { "char": 8, "short": 16, "int": 32 }
}
`
	if _, err := parseBytes([]byte(jsonStr)); err == nil {
		t.Error("'long' parameter can not be omitted")
	}
}
