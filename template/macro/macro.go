package macro

import (
	"errors"
	"regexp"
	"fmt"
	"strings"
)

type Macro struct {
	name string
	body string
	dummyArgs []string
}

func New(name string, body string, dummyArgs []string) (*Macro, error) {
	if name == "" {
		return nil, errors.New("mandatory paramter 'name' should not be empty")
	}

	return &Macro{name: name, body: body, dummyArgs: dummyArgs}, nil
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

	if len(macro.dummyArgs) != len(args) {
		err := fmt.Errorf("[%s]argument length is not be matched(expected=%d, got=%d)",
			macro.name, len(macro.dummyArgs), len(args))
		return "", err
	}

	bindings := make(map[string]string)
	for i, arg := range args {
		dummy := macro.dummyArgs[i]
		bindings[dummy] = arg;
	}

	matcheds := macroExpressionRegexp.FindAllStringSubmatch(macro.body, -1)
	if len(matcheds) == 0 {
		return "", nil
	}

	retval := macro.body
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
			if len(val.dummyArgs) == 0 {
				expanded = val.body
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
