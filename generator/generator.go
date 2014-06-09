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
	GlobalEnv map[string]*macro.Macro
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
	if err := os.Mkdir(generator.Config.TestDir, 0755); err != nil {
		return err
	}

	//	for _, template := range templates {
	//		// generate template file
	//	}

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

	generator.setPredefinedMacros()

	if err := generator.generateTestSuite(templates); err != nil {
		return err
	}

	return nil
}

func (generator *Generator) setPredefinedMacros() {
	size := generator.Config.Size

	generator.registerTypeMacros("char", size.Char)
	generator.registerTypeMacros("short", size.Short)
	generator.registerTypeMacros("int", size.Int)
	generator.registerTypeMacros("long", size.Long)
	// should implement pointer type and 'float' and 'double'
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
	return fmt.Sprintf("%g%s", math.Pow(2, width)-1, suffix)
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

func (generator *Generator) registerTypeMacros(typeName string, bitWidth int) {
	complement := generator.Config.Complement

	signedMin := signedMinValue(typeName, bitWidth, complement)
	signedMax := signedMaxValue(typeName, bitWidth)
	unsignedMax := unsignedMaxValue(typeName, bitWidth)

	signedMinName := macroTypeName(typeName, SIGNED_TYPE, MIN_VALUE)
	signedMaxName := macroTypeName(typeName, SIGNED_TYPE, MAX_VALUE)
	unsignedMinName := macroTypeName(typeName, UNSIGNED_TYPE, MIN_VALUE)
	unsignedMaxName := macroTypeName(typeName, UNSIGNED_TYPE, MAX_VALUE)

	generator.GlobalEnv[signedMinName] = &macro.Macro{Name: signedMinName, Body: signedMin}
	generator.GlobalEnv[signedMaxName] = &macro.Macro{Name: signedMinName, Body: signedMax}
	generator.GlobalEnv[unsignedMinName] = &macro.Macro{Name: signedMinName, Body: "0"}
	generator.GlobalEnv[unsignedMaxName] = &macro.Macro{Name: signedMinName, Body: unsignedMax}
}
