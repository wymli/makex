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
	registerCobraCmds(rootCmd, makexfile.Cmds, makexfile.Udfs)
}

var (
	True  bool = true
	False bool = false
)

func (makexfile *Makexfile) FixUDF() {
	for i := range makexfile.Udfs {
		if makexfile.Udfs[i].Used == nil {
			makexfile.Udfs[i].Used = &True
		}
	}
}

// registerCobraCmds registers recursively
func registerCobraCmds(parentCmd *cobra.Command, userCmds []Cmd, udfs []UDF) {
	for _, userCmd := range userCmds {
		cobraCmd := buildCobraCmd(userCmd, udfs)
		parentCmd.AddCommand(cobraCmd)
		registerCobraCmds(cobraCmd, userCmd.Cmds, udfs)
	}
}

func buildCobraCmd(userCmd Cmd, udfs []UDF) *cobra.Command {
	return &cobra.Command{
		Use:     userCmd.Name,
		Aliases: userCmd.Aliases,
		Run: func(cmd *cobra.Command, args []string) {
			// 1. write cmds to a shell file
			f, err := os.CreateTemp(os.TempDir(), "makex-*.sh")
			if err != nil {
				log.Fatalf("failed to create tmp shell file, err: %v", err)
			}

			_ = f.Chmod(os.ModePerm)

			// write imports
			udfMap := map[string]UDF{}
			for _, udf := range udfs {
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

				// use cmd first if cmd is not empty
				switch {
				case udf.Cmd != "":
					safeWriteString(f, udf.Cmd+"\n")
				case udf.Load != "":
					source := fmt.Sprintf("source %s\n", udf.Load)
					safeWriteString(f, source)
				default:
					log.Fatal("udf '%s' in cmd '%s' must set 'cmd' or 'load'", udf.Name, userCmd.Name)
				}
			}

			// write commands
			safeWriteString(f, userCmd.Cmd)

			_ = f.Sync()

			if log.IsLevelEnabled(log.DebugLevel) {
				// show the whole shell file
				data, _ := os.ReadFile(f.Name())
				log.Debugf("shell file:\n%s", string(data))
			}

			log.Debugf("[run] temp shell file: %s", f.Name())

			// 2. exec shell
			shell := exec.Command("sh", f.Name())
			shell.Stdout = os.Stdout
			shell.Stderr = os.Stderr
			shell.Stdin = os.Stdin

			err = shell.Run()
			log.Debugf("shell exited with error=%v", err)
		},
	}
}

// func WriteImports(w io.Writer, udf []UDF, imports []string) error {
// 	for _, shell := range imports {
// 	}
// 	return nil
// }

func safeWriteString(f *os.File, str string) {
	n, err := f.WriteString(str)
	if err != nil || n != len(str) {
		log.Fatalf("failed to write tmp shell file, writed %d bytes of %d bytes, err: %v", n, len(str), err)
	}
}
