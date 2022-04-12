package parser

type Makexfile struct {
	Interpreter string `yaml:"interpreter,omitempty"`
	Udfs        []UDF  `yaml:"udfs,omitempty"`
	Cmds        []Cmd  `yaml:"cmds,omitempty"`
}

type UDF struct {
	Name   string `yaml:"name,omitempty"`
	Prompt string `yaml:"prompt,omitempty"`
	Cmd    string `yaml:"cmd,omitempty"`
	Load   string `yaml:"load,omitempty"`
	Used   *bool  `yaml:"used,omitempty"`
}

type Cmd struct {
	Name    string   `yaml:"name,omitempty"`
	Aliases []string `yaml:"aliases,omitempty"`
	Imports []string `yaml:"imports,omitempty"`
	Cmd     string   `yaml:"cmd,omitempty"`
	Cmds    []Cmd    `yaml:"cmds,omitempty"`
}
