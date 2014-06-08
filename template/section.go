package template

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type sectionFunc func(parser *Parser, arg string, content string) error

var dispatchTable map[string]sectionFunc

func init() {
	dispatchTable = make(map[string]sectionFunc)
	dispatchTable["def"] = parseDefSection
	dispatchTable["dir"] = parseDirSection
	dispatchTable["include"] = parseIncludeSection
	dispatchTable["comment"] = parseCommentSection
}

func parseDefSection(parser *Parser, arg string, content string) error {
	return nil
}

var emptyLineRegexp = regexp.MustCompile(`^\s*$`)
var commentStart = regexp.MustCompile(`^@comment`)
var commentEnd = regexp.MustCompile(`^@comment_`)

func skipCommentSection(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		line := scanner.Text()
		if commentEnd.MatchString(line) {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

var macroCallRegexp = regexp.MustCompile(`^([^(]+)\((.*)\)$`)
var macroArgRegexp = regexp.MustCompile(`(?:(?:\[([^]])+\])|(\w+))\s*,?`)

type macroCall struct {
	name      string
	arguments []string
}

func parseMacroString(macroStr string) (*macroCall, error) {
	matched := macroCallRegexp.FindStringSubmatch(macroStr)
	if matched == nil {
		return nil, fmt.Errorf("Invalid macro: '%s'", macroStr)
	}

	name := matched[1]
	arg := matched[2]

	arguments := make([]string, 0)

	matcheds := macroArgRegexp.FindAllStringSubmatch(arg, 0)
	for _, matched := range matcheds {
		if matched[1] != "" { // argument in brackets
			arg := strings.Trim(matched[1], " \t\r\n")
			arguments = append(arguments, arg)
		} else if matched[2] != "" {
			arg := strings.Trim(matched[2], " \t\r\n")
			arguments = append(arguments, arg)
		}
	}

	macro := &macroCall{
		name:      name,
		arguments: arguments,
	}

	return macro, nil
}

// TODO should test this regexp
var fileSectionRegexp = regexp.MustCompile(`(?sm)@file\s+(\S+)\s+(\$[^(]+\(.*\))\s+(?:@ok\s+(\d+)\s+@ok_\s+)?\@file_`)

func processDirSection(content string) error {
	sr := strings.NewReader(content)
	scanner := bufio.NewScanner(sr)

	for scanner.Scan() {
		line := scanner.Text()
		if emptyLineRegexp.MatchString(line) {
			continue
		}

		if matched := fileSectionRegexp.FindStringSubmatch(line); matched != nil {
			// $1=filename, $2=macro(args), $3=oknum
			_, err := parseMacroString(matched[2]) // XXX
			if err != nil {
				return err
			}

		} else if commentStart.MatchString(line) {
			if commentEnd.MatchString(line) {
				continue
			}

			if err := skipCommentSection(scanner); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseDirSection(parser *Parser, arg string, content string) error {
	dir := strings.Trim(arg, " \t\n\r")

	dirPath := filepath.Join(parser.outputDirectory, dir)
	if _, err := os.Stat(dirPath); os.IsExist(err) {
		if err := os.RemoveAll(dirPath); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(dirPath, 0777); err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(dirPath); err != nil {
		return err
	}
	defer os.Chdir(pwd)

	return nil
}

func parseIncludeSection(parser *Parser, arg string, content string) error {
	path := strings.Trim(content, " \t\n\r")

	var includedFile string
	for _, includePath := range parser.IncludePaths {
		copyed := path
		if !filepath.IsAbs(path) {
			copyed = filepath.Join(includePath, copyed)
		}

		if _, err := os.Stat(copyed); !os.IsNotExist(err) {
			includedFile = copyed
			break
		}
	}

	if includedFile == "" {
		return fmt.Errorf("'%s' is not found", path)
	}

	file, err := os.Open(includedFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return parser.parseTemplate(file)
}

func parseCommentSection(parser *Parser, arg string, content string) error {
	// Do nothing
	return nil
}
