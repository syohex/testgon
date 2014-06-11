package template

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/syohex/testgon/template/macro"
	"path/filepath"
)

type Parser struct {
	IncludePaths     []string
	predefined       map[string]*macro.Macro
	filenameIndex    int
	templateEncoding string
	outputEncoding   string
	outputDirectory  string
}

func New(outputDir string, predefined map[string]*macro.Macro) *Parser {
	parser := &Parser{
		outputDirectory:  outputDir,
		predefined:       predefined,
		filenameIndex:    0,
		templateEncoding: "utf-8",
		outputEncoding:   "utf-8",
	}

	return parser
}

var startSection = regexp.MustCompile(`^@([^_\s]+)`)
var endSection = regexp.MustCompile(`^@([^_\s]+)_`)

func checkSyntax(file io.Reader) error {
	type section struct {
		name string
		line int
	}

	scanner := bufio.NewScanner(file)
	currentLine := 1

	sections := make([]*section, 0)
	currentSection := 0
	for scanner.Scan() {
		line := scanner.Text()

		if matched := endSection.FindStringSubmatch(line); matched != nil {
			section := matched[1]

			if currentSection == 0 {
				return fmt.Errorf("found only '%s' end at %d",
					section, currentLine)
			}

			lastSection := sections[currentSection-1].name
			if section != lastSection {
				return fmt.Errorf("missing end of '%s' section(at %d)",
					lastSection, sections[currentSection-1].line)
			}
			currentSection--
		} else if matched := startSection.FindStringSubmatch(line); matched != nil {
			sectionName := matched[1]

			re, err := regexp.Compile(`^@` + sectionName + `_\s*$`)
			if err != nil {
				return fmt.Errorf("can't create regexp object for '@%s'",
					sectionName)
			}
			if !re.MatchString(line) { // End section directive is not same line
				s := &section{
					name: sectionName,
					line: currentLine,
				}

				if currentSection < cap(sections) {
					sections[currentSection] = s
				} else {
					sections = append(sections, s)
				}
				currentSection++
			}
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if currentSection != 0 {
		messages := make([]string, 0)

		for i := 0; i < currentSection; i++ {
			msg := fmt.Sprintf("'%s' section doesn't have end of section(at %d)",
				sections[i].name, sections[i].line)
			messages = append(messages, msg)
		}

		return errors.New(strings.Join(messages, "\n"))
	}

	return nil
}

func (parser *Parser) parseTemplate(template io.Reader) error {
	scanner := bufio.NewScanner(template)

	for scanner.Scan() {
		line := scanner.Text()

		matched := startSection.FindStringSubmatch(line)
		if matched == nil {
			continue
		}

		section := matched[1]
		argument := matched[2]
		callback := dispatchTable[section]

		endRegexp, err := regexp.Compile(`^` + section + `_`)
		if err != nil {
			return errors.New("Can't create end regexp XXX")
		}

		lines := make([]string, 0)
		for {
			if endRegexp.MatchString(line) {
				break
			}
			lines = append(lines, line)
		}

		callback(parser, argument, strings.Join(lines, "\n"))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (parser Parser) Parse(template string) error {
	// Set directory in template file as default include path
	abs, err := filepath.Abs( filepath.Dir(template) )
	if err != nil {
		return err
	}
	parser.IncludePaths = append(parser.IncludePaths, abs)

	file, err := os.Open(template)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := checkSyntax(file); err != nil {
		return err
	}

	if err := parser.parseTemplate(file); err != nil {
		return err
	}

	// clean up default include path
	parser.IncludePaths = parser.IncludePaths[:len(parser.IncludePaths)-1]
	return nil
}
