# makex
It's a cmd-line tool like `make` and `task`, supporting nested options and alias using `cobra`.


## Shell
we use `sh`, not `bash` or other shell.

## Example

A `makefile` in `makex`, is named `makex.yaml`. We should place it in the root dir of your project.  

First, init.
```
cd 
```



```yaml
# user defined functions
udfs:
  - name: builtin # udf name
    prompt: color # prompt is just a kind of promption, meaning we have color file in built-in
    used: false   # we won't load it as udf, but we do have it in built-in
  - name: echo
    prompt: this is an example to load shell from multi-line text
    cmd: | # we will exec the cmd
      echo "hello world snippet_1"
      echo "hello world snippet_1"
  - name: echofile
    prompt: this is an example to load shell from file
    load: ./echofile.sh # we will source the shell file. relative path with makex.yaml

# running with `makex init, makex api init, makex api db init`
cmds:
  - name: init
    aliases: [create]
    imports: [echofile]
    cmd: |
      echo "init called"

  - name: api
    cmds:
      - name: init
        aliases: [create]
        imports: [color] # import will first search udf, then builtin.
        cmd: |
          ECHO_RED "api init called"

      - name: db
        aliases: [database]
        cmds:
          - name: init
            aliases: [create]
            cmd: |
              echo "api db init called"

```

```
liwm29@wymli-NB1:~/test$ makex init
echofile called
init called

liwm29@wymli-NB1:~/test$ makex api init 
api init called
```

## Feature
- local shell file load support
- built-in shell functions support

