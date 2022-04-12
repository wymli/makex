/*
Copyright © 2022 Li Weiming <liwm29@mail2.sysu.edu.cn>

*/
package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/wymli/makex/cmd"
	"github.com/wymli/makex/internal/parser"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	formatter := log.TextFormatter{
		EnvironmentOverrideColors: true,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
		QuoteEmptyFields:          true,
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			seps := strings.Split(frame.File, string(filepath.Separator))
			for i := len(seps); i < 3; i++ {
				seps = append(seps, "")
			}

			seps = seps[len(seps)-3:]

			fileName := filepath.Join(seps...) + ":" + strconv.Itoa(frame.Line)
			return "", fileName
		},
	}
	log.SetFormatter(&formatter)
}

func main() {
	// we don't process error here, just parse flags
	_ = cmd.RootCmd.ParseFlags(os.Args[1:])

	if makexfile, err := cmd.RootCmd.Flags().GetString("makexfile"); err == nil {
		if makexfile != "" {
			viper.Set("makexfile", makexfile)
		}
	}

	if v, err := cmd.RootCmd.Flags().GetBool("verbose"); err == nil {
		if v {
			log.SetLevel(log.TraceLevel)
		}
	}

	// if it is registered args, exec it directly; otherwise we register user-defined args and re-exec again.
	// args := cmd.RootCmd.Flags().Args()
	// // if we don't have args, we register user-defined args and then show help info
	// if _, _, err := cmd.RootCmd.Find(args); err == nil && len(args) != 0 {
	// 	log.Debugf("[run] exec builtin cmd, args: %v", args)
	// 	cmd.Execute()
	// 	return
	// }

	// 1. get makexfilename name
	makexfilename := viper.GetString("makexfile")

	log.Debugf("use makexfile: %s", makexfilename)

	// 2. read makexfile
	makexfile, err := parser.ReadMakexfile(makexfilename)
	if err != nil {
		log.Fatal(err)
	}

	// 3. register makexfile cmds to cobra cmds
	log.Debugf("[register] registering commands to cobra")
	makexfile.RegisterCmds(cmd.RootCmd)

	// 4. run cobra cmds executor
	log.Debugf("[execute] executing cobra")
	cmd.Execute()
}
