package generator

import (
	"github.com/syohex/testgon/config"
	"github.com/syohex/testgon/template/macro"
	"errors"
	"path/filepath"
	"fmt"
	"regexp"
	"os"
	"math"
)

type Param struct {
	File string
	Help boolean
	IntOnly boolean
	FloatOnly boolean
}

type Generator struct {
	Config config.Config
	Help boolean
	IntOnly boolean
	FloatOnly boolean
}

func New(param Config) (Generator, error){
	conf, err := config.Parse(param.File)
	if err != nil {
		return nil, err
	}

	generator := &Generator{
		Config: conf,
		Help: param.Help,
		IntOnly: param.IntOnly,
		FloatOnly: param.FloatOnly
	}

	return generator, nil
}

func expandPatterns(patterns []string) ([]string, error) {
	files = []string
	for _, pattern := range patterns {
		expands, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("Can't expand '%s'", pattern)
		}

		files = append(files, ...expands)
	}

	return files
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

func (generator *Generator)generateTestSuite(templates []string) {
	if err := os.Mkdir(generator.Config.TestDir); err != nil {
		return err
	}

	for _, template := range templates {
		// generate template file
	}
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

	generator.generaateTestSuite(templates)
}

var spaceRe = regexp.MustCompile(`\s+`)

func setPredefinedMacros(conf config.Config) {
	typeSuffix := map[string]string {
		"long": "L",
		"long long": "LL",
		"float": "F",
		"long double": "L",
	}

	predefined := make(map[string]*macro.Macro)
	for typeName, size := range conf.Size {
		typeName = spaceRe.ReplaceAllString(typeName, "")
	}
}

type Limit struct {
	signedCharMin int64
	signedCharMax int64
	unsighedCharMin uint64
	unsignedCharMax uint64

	singedShortMin int64
	singedShortMax int64
	unsignedShortMin int64
	unsignedShortMax int64

	signedIntMin int64
	signedIntMax int64
	unsignedIntMin int64
	unsignedIntMax int64

	signedLongMin int64
	singedLongMax int64
	unsignedLongMin int64
	unsignedLongMax int64

	signedLongLongMin int64
	signedLongLongMax int64
	unsignedLongLongMin int64
	unsignedLongLongMax int64
}

func getTypeLimit(size config.Size, complement int) *Limit {
	limit := &Limit{}

	// Use string and fmt.Sprintf for representing limit number
	for typeName, bitWidth := range size {
		switch typeName {
		case "char":
			limit.signedCharMin = -(2 ** (bitWidth - 1))
			limit.signedCharMax = (2 ** (bitWidth - 1)) - 1
			limit.unsignedCharMax = (2 ** bitWidth) - 1
		case "short":
			limit.signedShortMin = -(2 ** (bitWidth - 1))
			limit.signedShortMax = (2 ** (bitWidth - 1)) - 1
			limit.unsignedShortMax = (2 ** bitWidth) - 1
		case "int":
			limit.signedIntMin = -(2 ** (bitWidth - 1))
			limit.signedIntMax = (2 ** (bitWidth - 1)) - 1
			limit.unsignedIntMax = (2 ** bitWidth) - 1
		case "long":
			limit.signedLongMin = -(2 ** (bitWidth - 1))
			limit.signedLongMax = (2 ** (bitWidth - 1)) - 1
			limit.unsignedLongMax = (2 ** bitWidth) - 1
		case "longlong":
			limit.signedLongLongMin = -(2 ** (bitWidth - 1))
			limit.signedLongLongMax = (2 ** (bitWidth - 1)) - 1
			limit.unsignedLongLongMax = (2 ** bitWidth) - 1
		}
	}

	if complement == 1 {
		limit.signedCharMin = -(2 ** (bitWidth - 1))
		limit.signedShortMin = -(2 ** (bitWidth - 1))
		limit.signedIntMin = -(2 ** (bitWidth - 1))
		limit.signedLongMin = -(2 ** (bitWidth - 1))
		limit.signedLongLongMin = -(2 ** (bitWidth - 1))
	}

	return limit
}
