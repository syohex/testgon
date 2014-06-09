package generator

import (
	"github.com/syohex/testgon/config"
	"github.com/syohex/testgon/template/macro"
	"errors"
	"path/filepath"
	"fmt"
	"regexp"
	"os"
	"strings"
	"math"
)

type Param struct {
	File string
	Help bool
	IntOnly bool
	FloatOnly bool
}

type Generator struct {
	Config *config.Config
	Help bool
	IntOnly bool
	FloatOnly bool
}

func New(param Param) (*Generator, error){
	conf, err := config.Parse(param.File)
	if err != nil {
		return nil, err
	}

	generator := &Generator{
		Config: conf,
		Help: param.Help,
		IntOnly: param.IntOnly,
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

func (generator *Generator)generateTestSuite(templates []string) error {
	if err := os.Mkdir(generator.Config.TestDir, 0755); err != nil {
		return err
	}

//	for _, template := range templates {
//		// generate template file
//	}

	return nil
}

func (generator *Generator)Run(patterns []string) error {
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

	// TODO set macros

	if err := generator.generateTestSuite(templates); err != nil {
		return err
	}

	return nil
}

var spaceRe = regexp.MustCompile(`\s+`)

func setPredefinedMacros(conf config.Config) map[string]*macro.Macro {
	predefined := make(map[string]*macro.Macro)
	registerTypeMacros(predefined, "char", conf.Size.Char, conf.Complement)
	registerTypeMacros(predefined, "short", conf.Size.Short, conf.Complement)
	registerTypeMacros(predefined, "int", conf.Size.Int, conf.Complement)
	registerTypeMacros(predefined, "long", conf.Size.Long, conf.Complement)
	// should implement pointer type and 'float' and 'double'

	return predefined
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
	suffix := typeSuffix(typeName);
	width := float64(bitWidth)
	return fmt.Sprintf("%g%s", math.Pow(2, (width - 1)) - 1, suffix)
}

func signedMinValue(typeName string, bitWidth int, complement int) string {
	suffix := typeSuffix(typeName);

	width := float64(bitWidth)
	if complement == 2 {
		return fmt.Sprintf("%g%s", -math.Pow(2, (width - 1)), suffix)
	} else {
		return fmt.Sprintf("%g%s", -math.Pow(2 , (width - 1)) + 1, suffix)
	}
}

func unsignedMaxValue(typeName string, bitWidth int) string {
	suffix := typeSuffix(typeName);
	width := float64(bitWidth)
	return fmt.Sprintf("%g%s", math.Pow(2, width) - 1, suffix)
}

// use Enum type instead of bool type ???
func macroTypeName(typeName string, unsigned bool, min bool) string {
	var unsignedPrefix string
	if unsigned {
		unsignedPrefix = "U"
	} else {
		unsignedPrefix = ""
	}

	var suffix string
	if min {
		suffix = "MIN"
	} else {
		suffix = "MAX"
	}

	return fmt.Sprintf("%s%s%s", unsignedPrefix, strings.ToUpper(typeName), suffix)
}

func registerTypeMacros(
	predefined map[string]*macro.Macro,
	typeName string,
	bitWidth int,
	complement int,
) {
	signedMin := signedMinValue(typeName, bitWidth, complement)
	signedMax := signedMaxValue(typeName, bitWidth)
	unsignedMax := unsignedMaxValue(typeName, bitWidth)

	signedMinName := macroTypeName(typeName, false, true)
	signedMaxName := macroTypeName(typeName, false, false)
	unsignedMinName := macroTypeName(typeName, true, true)
	unsignedMaxName := macroTypeName(typeName, true, false)

	predefined[signedMinName] = &macro.Macro{Name: signedMinName, Body: signedMin}
	predefined[signedMaxName] = &macro.Macro{Name: signedMinName, Body: signedMax}
	predefined[unsignedMinName] = &macro.Macro{Name: signedMinName, Body: "0"}
	predefined[unsignedMaxName] = &macro.Macro{Name: signedMinName, Body: unsignedMax}
}
