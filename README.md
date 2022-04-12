# makex
It's a cmd-line tool like `make` and `task`, supporting nested options and alias using `cobra`.  
With `makex`, we can easily execute nested commands, like `makex rpc build`, `makex rpc build pb`

## Usage
you can run `makex template init` to generate makexfile(`makex.yaml`) in `$pwd`, then you can edit your own makex.yaml and run with `makex`, just like `make in cobra style`

For Example, if you have `init` cmd in your own `makex.yaml`, you can run `makex init` in the same dir of your `makex.yaml` to execute init commands defined in `makex.yaml`. You can type `makex help init (makex init -h, makex init --help)` to see help information. If `init` is an empty cmd, `makex init` will also print help info.

> Normally, you can just type `makex, makex help, makex -h, makex --help` to get help info.

> More cli usage, you can ask help for `cobra doc`.

## Example

A `makefile` in `makex`, is named `makex.yaml`. We should place it in the root dir of your project.  

[![asciicast](https://asciinema.org/a/486509.svg)](https://asciinema.org/a/486509)


``` yaml
interpreter: sh  # see https://pubs.opengroup.org/onlinepubs/9699919799/utilities/contents.html

# user defined functions
udfs:
  - name: builtin   # this entry is just for promption, meaning we have `color` in builtin functions
    prompt: color
    used: false
  - name: genpb
    cmd: |
      genpb(){
        protoc -I $dir \
        --go_out $out_dir --go_opt paths=source_relative \
        --go-grpc_out $out_dir --go-grpc_opt paths=source_relative \
        *.proto
      }

# running with `makex init, makex tidy, makex userrpc build, makex userrpc build pb`
# run `makex help` will give you a list
cmds:
  - name: init
    cmd: |
      go mod init github.com/wymli/makex_example
  - name: tidy
    cmd: |
      go mod tidy
  - name: userrpc
    aliases: []
    usage: userrpc is abould user center rpc
    imports: []
    cmds:
      - name: build
        aliases: [gen]
        cmd: |
          cd user_rpc
          go build
          # go build user_rpc/main.go -o bin/user_rpc
        cmds:
          - name: pb
            imports: [genpb]
            cmd: |
              dir=.
              out_dir=.
              cd user_rpc/proto
              genpb
      - name: run
        cmd: |
          ./user_rpc/user_rpc
```

## Exec
We organize all shell commands into one big temp shell file.  
We first process imports in `cmd`
- if the udf of imports has `cmd` field, we will copy udf.cmd to the shell file
- if the udf of imports has `load` field, we will use `. ${udf.load}` to `source`

Then we copy `cmd.cmd` to shell file and run it using a shell interpreter(default `sh`).

## Makexfile(makex.yaml) Schema
```
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
  Usage   string   `yaml:"usage,omitempty"`
	Imports []string `yaml:"imports,omitempty"`
	Cmd     string   `yaml:"cmd,omitempty"`
	Cmds    []Cmd    `yaml:"cmds,omitempty"`
}
```

### UDF
UDF can be seen as a kind of code snippets.

- Name： used when imported in cmd (Cmd.Imports)
- Prompt: useful information, like exported funtion name list
- Cmd/Load: the payload, if cmd is not empty, we will exec cmd, otherwise will load shell file. cmd should be shell commands, and load should be shell file(relative to the makex.yaml)
- Used: use or not

### Cmd
Usage is just like cobra.

- Name: cmd name
- Aliases: cmd alias
- Usage: info showed in help usage
- Imports: using udf
- Cmd: command to execute
- Cmds: sub-commands



## Shell
we use [`sh`](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/contents.html), not `bash` or other shell.  
> you can change `interpreter` in `makex.yaml` easily, but normally builtin shell function is coded using `sh`

> a marked difference between different shell is that when showing color in echo, `sh` is just `echo` without -e, while `bash` needs `echo -e`


## Config
you can configure you owm template on makexfile(makex.yaml), which is located at `$HOME/.makex/makex_config.yaml`.
- do `cat ~/.makex/makex_config.yaml` for detail

you can store your own commonly used udf(code snippets) as builtin functions at `$HOME/.makex/shell/`
- do `ls ~/.makex/shell/` for detail
- each udf is organized as a file
  - filename withoud ext is its import name
  - a file can contains as many functions as you want
  - we will load the whole file before `exec cmd`

## Debug

you can run with `-v` to show debug logs. The contents of the assembled shell file will also be displayed.