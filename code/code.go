package code

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//go:embed code
var CodeFS embed.FS

// RegisterCmds registers recursively
func RegisterCmds(parentCmd *cobra.Command) {
	cmdMap := map[string]*cobra.Command{
		"code": parentCmd,
	}

	fs.WalkDir(CodeFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		names := strings.Split(path, string(os.PathSeparator))
		log.Debug("path=", path, "d.name=", d.Name(), "names=", names)

		for i, name := range names[:len(names)-1] {
			if _, ok := cmdMap[name]; !ok {
				cmd := &cobra.Command{Use: name}
				cmdMap[name] = cmd
				cmdMap[names[i-1]].AddCommand(cmd)
			}
		}

		realCmd := buildCobraCmd(path, d.Name())
		cmdMap[names[len(names)-2]].AddCommand(realCmd)

		return nil
	})
}

func buildCobraCmd(path, name string) *cobra.Command {
	nameSplited := strings.Split(name, ".")

	return &cobra.Command{
		Use: nameSplited[0],
		Run: func(cmd *cobra.Command, args []string) {
			// 1.copy the file to $pwd
			cwd, _ := os.Getwd()

			filepath.Join(cwd, name)

			data, err := CodeFS.ReadFile(path)
			if err != nil {
				log.Fatalf("failed to open embedfs file, err: %v", err)
			}

			err = os.WriteFile(filepath.Join(cwd, name), data, os.ModePerm)
			if err != nil {
				log.Fatalf("faild to write file to cwd '%s', err: %v", cwd, err)
			}
		},
	}
}
