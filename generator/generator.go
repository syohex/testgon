package generator

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/syohex/testgon/config"
	"github.com/syohex/testgon/template"
	"github.com/syohex/testgon/template/macro"
)

type Param struct {
	File      string
	Help      bool
	IntOnly   bool
	FloatOnly bool
}

type Generator struct {
	Config    *config.Config
	Help      bool
	IntOnly   bool
	FloatOnly bool
}

func New(param Param) (*Generator, error) {
	conf, err := config.Parse(param.File)
	if err != nil {
		return nil, err
	}

	generator := &Generator{
		Config:    conf,
		Help:      param.Help,
		IntOnly:   param.IntOnly,
		FloatOnly: param.FloatOnly,
	}

	return generator, nil
}

func expandPatterns(patterns []string) ([]string, error) {
	files := make([]string, 0)
	for _, pattern := range patterns {
		expands, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("Can't expand '%s'", pattern)
		}

		files = append(files, expands...)
	}

	return files, nil
}

var templateFileSuffix = regexp.MustCompile(`\.tt$`)

func checkTemplateFileName(templates []string) error {
	for _, template := range templates {
		if !templateFileSuffix.MatchString(template) {
			return errors.New("Suffix of template file should be '.tt'")
		}
	}

	return nil
}

func (generator *Generator) generateTestSuite(templates []string) error {
	env := generator.setPredefinedMacros()

	outputDir := generator.Config.TestDir
	if err := os.Mkdir(outputDir, 0755); err != nil {
		return err
	}

	parser := template.NewParser(outputDir, env)
	for _, template := range templates {
		if err := parser.Parse(template); err != nil {
			return err
		}
	}

	return nil
}

func (generator *Generator) Run(patterns []string) error {
	if patterns == nil {
		return errors.New("Templete files are not specified")
	}

	templates, err := expandPatterns(patterns)
	if err != nil {
		return err
	}

	if err := checkTemplateFileName(templates); err != nil {
		return err
	}

	if err := generator.generateTestSuite(templates); err != nil {
		return err
	}

	return nil
}

func (generator *Generator) setPredefinedMacros() map[string]*macro.Macro {
	size := generator.Config.Size
	complement := generator.Config.Complement

	env := make(map[string]*macro.Macro)
	registerIntTypeMacro(env, "char", size.Char, complement)
	registerIntTypeMacro(env, "short", size.Short, complement)
	registerIntTypeMacro(env, "int", size.Int, complement)
	registerIntTypeMacro(env, "long", size.Long, complement)
	// should implement pointer type and 'float' and 'double'

	return env
}

func typeSuffix(typeName string) string {
	switch typeName {
	case "long":
		return "L"
	case "long long":
		return "LL"
	case "float":
		return "F"
	case "long double":
		return "L"
	default:
		return ""
	}
}

func signedMaxValue(typeName string, bitWidth int) string {
	suffix := typeSuffix(typeName)
	width := float64(bitWidth)
	return fmt.Sprintf("%.0f%s", math.Pow(2, (width-1))-1, suffix)
}

func signedMinValue(typeName string, bitWidth int, complement int) string {
	suffix := typeSuffix(typeName)

	width := float64(bitWidth)
	if complement == 2 {
		return fmt.Sprintf("%.0f%s", -math.Pow(2, (width-1)), suffix)
	} else {
		return fmt.Sprintf("%.0f%s", -math.Pow(2, (width-1))+1, suffix)
	}
}

func unsignedMaxValue(typeName string, bitWidth int) string {
	suffix := typeSuffix(typeName)
	width := float64(bitWidth)
	return fmt.Sprintf("%.0f%s", math.Pow(2, width)-1, suffix)
}

const (
	SIGNED_TYPE   = 0
	UNSIGNED_TYPE = 1
)

const (
	MIN_VALUE = 0
	MAX_VALUE = 1
)

func macroTypeName(typeName string, signed int, min int) string {
	var unsignedPrefix string
	if signed == UNSIGNED_TYPE {
		unsignedPrefix = "U"
	} else {
		unsignedPrefix = ""
	}

	var suffix string
	if min == MIN_VALUE {
		suffix = "MIN"
	} else {
		suffix = "MAX"
	}

	return fmt.Sprintf("%s%s%s", unsignedPrefix, strings.ToUpper(typeName), suffix)
}

func registerIntTypeMacro(
	env map[string]*macro.Macro,
	typeName string,
	bitWidth int,
	complement int,
) {
	signedMin := signedMinValue(typeName, bitWidth, complement)
	signedMax := signedMaxValue(typeName, bitWidth)
	unsignedMax := unsignedMaxValue(typeName, bitWidth)

	signedMinName := macroTypeName(typeName, SIGNED_TYPE, MIN_VALUE)
	signedMaxName := macroTypeName(typeName, SIGNED_TYPE, MAX_VALUE)
	unsignedMinName := macroTypeName(typeName, UNSIGNED_TYPE, MIN_VALUE)
	unsignedMaxName := macroTypeName(typeName, UNSIGNED_TYPE, MAX_VALUE)

	env[signedMinName] = &macro.Macro{Name: signedMinName, Body: signedMin}
	env[signedMaxName] = &macro.Macro{Name: signedMinName, Body: signedMax}
	env[unsignedMinName] = &macro.Macro{Name: signedMinName, Body: "0"}
	env[unsignedMaxName] = &macro.Macro{Name: signedMinName, Body: unsignedMax}
}
