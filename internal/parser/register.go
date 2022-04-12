package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/wymli/makex/internal/config"
	"github.com/wymli/makex/shell"
)

func ReadMakexfile(filename string) (*Makexfile, error) {
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("makexfile '%s' not exists, try 'makex template init' first", filename)
	} else if err != nil {
		return nil, fmt.Errorf("failed to read makexfile '%s', err: %v", filename, err)
	}

	makexfile := Makexfile{}

	err = yaml.Unmarshal(data, &makexfile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall file %s to go struct, err: %v", filename, err)
	}

	makexfile.FixUDF()

	out, _ := json.Marshal(makexfile)
	log.Tracef("makexfile in memory: \n%s", string(out))

	return &makexfile, nil
}

func (makexfile *Makexfile) RegisterCmds(rootCmd *cobra.Command) {
	registerCobraCmds(rootCmd, makexfile.Cmds, makexfile)
}

var (
	True  bool = true
	False bool = false
)

func (makexfile *Makexfile) FixUDF() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd, err: %v", err)
	}

	for i := range makexfile.Udfs {
		if makexfile.Udfs[i].Used == nil {
			makexfile.Udfs[i].Used = &True
		}
		if makexfile.Udfs[i].Load != "" {
			makexfile.Udfs[i].Load = filepath.Join(cwd, makexfile.Udfs[i].Load)
		}
	}
}

// registerCobraCmds registers recursively
func registerCobraCmds(parentCmd *cobra.Command, cmds []Cmd, makexfile *Makexfile) {
	for _, cmd := range cmds {
		cobraCmd := buildCobraCmd(cmd, makexfile)
		parentCmd.AddCommand(cobraCmd)
		registerCobraCmds(cobraCmd, cmd.Cmds, makexfile)
	}
}

func buildCobraCmd(userCmd Cmd, makexfile *Makexfile) *cobra.Command {
	return &cobra.Command{
		Use:     userCmd.Name,
		Aliases: userCmd.Aliases,
		Short:   userCmd.Usage,
		Run: func(cmd *cobra.Command, args []string) {
			if userCmd.Cmd == "" {
				// if cmd is empty, we just print usage
				_ = cmd.Usage()
				return
			}
			// 1. create tmp shell file
			f, err := os.CreateTemp(os.TempDir(), "makex-*.sh")
			if err != nil {
				log.Fatalf("failed to create tmp shell file, err: %v", err)
			}

			_ = f.Chmod(os.ModePerm)

			// 2. write imports udf
			udfMap := map[string]UDF{}
			for _, udf := range makexfile.Udfs {
				if *udf.Used {
					udfMap[udf.Name] = udf
				}
			}

			builtinMap := map[string]UDF{}
			shells, err := shell.ShellFS.ReadDir(".")
			if err != nil {
				log.Fatalf("failed to read dir files from embed shell fs, err: %v", err)
			}

			for _, shell := range shells {
				name := strings.Split(shell.Name(), ".")[0]
				builtinMap[name] = UDF{
					Name: name,
					Load: filepath.Join(config.SHELL_DIR, shell.Name()),
					Used: &True,
				}
			}

			if log.IsLevelEnabled(log.DebugLevel) {
				udfNames := []string{}
				for udfName := range udfMap {
					udfNames = append(udfNames, udfName)
				}

				udfLen := len(udfNames)

				for builtinName := range builtinMap {
					udfNames = append(udfNames, builtinName)
				}

				log.Debugf("[build cmd] udf: %v builtin: %v", udfNames[:udfLen], udfNames[udfLen:])
			}

			log.Debugf("udf imports: %v", userCmd.Imports)

			for _, importt := range userCmd.Imports {
				var udf UDF
				var ok bool

				udf, ok = udfMap[importt]
				if !ok {
					udf, ok = builtinMap[importt]
					if !ok {
						log.Fatalf("unknown imports '%s' in cmd '%s'", importt, userCmd.Name)
					}
				}

				// use cmd first if cmd is not empty, then load
				switch {
				case udf.Cmd != "":
					safeWriteString(f, udf.Cmd+"\n")
				case udf.Load != "":
					source := fmt.Sprintf(". %s\n", udf.Load)
					safeWriteString(f, source)
				default:
					log.Fatal("udf '%s' import in cmd '%s', must set 'cmd' or 'load'", udf.Name, userCmd.Name)
				}
			}

			// 3. write commands
			safeWriteString(f, userCmd.Cmd)

			_ = f.Sync()

			if log.IsLevelEnabled(log.DebugLevel) {
				// show the whole shell file
				data, _ := os.ReadFile(f.Name())
				log.Debugf("shell file:\n%s", string(data))
			}

			log.Debugf("[run] temp shell file: %s", f.Name())

			// 4. exec shell file
			log.Debugf("[run] using shell: %s", makexfile.Interpreter)
			shell := exec.Command(makexfile.Interpreter, f.Name())
			shell.Stdout = os.Stdout
			shell.Stderr = os.Stderr
			shell.Stdin = os.Stdin

			err = shell.Run()
			log.Debugf("shell exited with error=%v", err)
		},
	}
}

func safeWriteString(f *os.File, str string) {
	n, err := f.WriteString(str)
	if err != nil || n != len(str) {
		log.Fatalf("failed to write tmp shell file, writed %d bytes of %d bytes, err: %v", n, len(str), err)
	}
}
