package macro

import (
	"regexp"
	"fmt"
	"strings"
)

type Macro struct {
	Name string
	Body string
	DummyArgs []string
}

var ignored_expressions = []string{"$Id"}

func isIgnoredExpression(expression string) bool {
	for _, ignored := range ignored_expressions {
		if expression == ignored {
			return true
		}
	}

	return false
}

var dollarRegex = regexp.MustCompile(`\$`)

func escapeRegexpString(regexpStr string) string {
	return fmt.Sprintf(`\Q%s\E`, regexpStr)
}

// macroExpressionRegexp matches expression in macro. (ex $foo(bar, baz))
// $1='macro name', $2='macro arguments'
var macroExpressionRegexp = regexp.MustCompile(`(\$[a-zA-Z0-9]+)(?:\(([^)]*)\))?`)

var spaceRegexp = regexp.MustCompile(`[\s\r\n]+`)
func removeAllSpaces(arg string) string {
	return spaceRegexp.ReplaceAllString(arg, "")
}

func (macro *Macro)Evaluate(args []string, env map[string]*Macro) (string, error) {

	if len(macro.DummyArgs) != len(args) {
		err := fmt.Errorf("[%s]argument length is not be matched(expected=%d, got=%d)",
			macro.Name, len(macro.DummyArgs), len(args))
		return "", err
	}

	bindings := make(map[string]string)
	for i, arg := range args {
		dummy := macro.DummyArgs[i]
		bindings[dummy] = arg;
	}

	matcheds := macroExpressionRegexp.FindAllStringSubmatch(macro.Body, -1)
	if len(matcheds) == 0 {
		return "", nil
	}

	retval := macro.Body
	for _, matched := range matcheds {
		name := matched[1]

		var expanded string
		if val, ok := bindings[name]; ok {
			expanded = val
		} else if val, ok := env[name]; ok {
			var args []string
			if matched[2] == "" {
				args = nil
			} else {
				argsStr := removeAllSpaces(matched[2])
				args = strings.Split(argsStr, ",")
			}

			var err error
			if len(val.DummyArgs) == 0 {
				expanded = val.Body
			} else {
				expanded, err = val.Evaluate(args, env)
				if err != nil {
					return "", err
				}
			}
		} else {
			if !isIgnoredExpression(name) {
				return "", fmt.Errorf("'%s' is not defined macro", name)
			}
			continue
		}

		replacedRegexp, err := regexp.Compile(regexp.QuoteMeta(matched[0]))
		if err != nil {
			return "", err
		}

		retval = replacedRegexp.ReplaceAllString(retval, expanded)
	}

	return retval, nil
}
