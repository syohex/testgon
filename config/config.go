package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"os"
	"io/ioutil"
)

// Config describes Testgen configuration
type Config struct {
	Compiler    string   `json:"compiler"`
	Simulator   string   `json:"simulator"`
	CFlags      []string `json:"c_flags"`
	LDFlags     []string `json:"ld_flags"`
	Options     []string `json:"options"`
	TestDir     string   `json:"testdir"`
	CompileOnly bool     `json:"compile_only"`

	Size integerTypeSize `json:"size"`

	Timeout   int  `json:"timeout"`
	Parallels int  `json:"parallels"`
	Color     bool `json:color`

	Lang            string `json:"lang"`
	Expect          string `json:"expect"`
	Complement      int    `json:"complement"`
	OutputOption    string `json:"output_option"`
	OptionSeparator string `json:"option_separator"`

	Temp *bool `json:"has_printf"`
	HasPrintf bool
}

type integerTypeSize struct {
	Char    int `json:"char"`
	Short   int `json:"short"`
	Int     int `json:"int"`
	Long    int `json:"long"`
	Pointer int `json:"pointer"`
}

func (conf *Config) checkMandatoryParameters() error {
	if conf.Compiler == "" {
		return errors.New("Not specified 'compiler' parameter")
	}

	if conf.TestDir == "" {
		return errors.New("Not specified 'testdir' parameter")
	}

	return nil
}

func (conf *Config) lookCommands() error {
	if _, err := exec.LookPath(conf.Compiler); err != nil {
		return fmt.Errorf("Compiler '%s' is not found in PATH")
	}

	if conf.Simulator != "" {
		if _, err := exec.LookPath(conf.Simulator); err != nil {
			return fmt.Errorf("Simulator '%s' is not found in PATH", conf.Simulator)
		}
	}

	return nil
}

func (size *integerTypeSize) checkSizeParameter() error {
	if size.Char == 0 {
		return errors.New("'char' in 'size' is not specified")
	}

	if size.Short == 0 {
		return errors.New("'short' in 'size' is not specified")
	}

	if size.Int == 0 {
		return errors.New("'int' in 'size' is not specified")
	}

	if size.Long == 0 {
		return errors.New("'long' in 'size' is not specified")
	}

	return nil
}

func (conf *Config) validate() error {
	if err := conf.checkMandatoryParameters(); err != nil {
		return err
	}

	if err := conf.lookCommands(); err != nil {
		return err
	}

	if err := conf.Size.checkSizeParameter(); err != nil {
		return err
	}

	return nil
}

func (conf *Config) setDefaultValue() {
	if conf.Lang == "" {
		conf.Lang = "c"
	}

	if conf.Parallels == 0 {
		conf.Parallels = 1
	}

	if conf.Timeout == 0 {
		conf.Timeout = 10
	}

	if conf.Expect == "" {
		conf.Expect = "@OK@"
	}

	if conf.Complement == 0 {
		conf.Complement = 2
	}

	if conf.OutputOption == "" {
		conf.OutputOption = "-o"
	}

	if conf.OptionSeparator == "" {
		conf.OptionSeparator = " "
	}

	if conf.Temp == nil {
		conf.HasPrintf = true // default value
	} else {
		conf.HasPrintf = *conf.Temp
	}
}

func Parse(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return parseBytes(bytes)
}

// Parse configuration file
func parseBytes(jsonBytes []byte) (*Config, error) {
	config := new(Config)

	if err := json.Unmarshal(jsonBytes, config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	config.setDefaultValue()

	return config, nil
}
