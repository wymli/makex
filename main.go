/*
Copyright © 2022 Li Weiming <liwm29@mail2.sysu.edu.cn>

*/
package main

import (
	"os"

	"github.com/wymli/makex/cmd"
	"github.com/wymli/makex/cmd/template"
	"github.com/wymli/makex/code"
	"github.com/wymli/makex/internal/config"
	"github.com/wymli/makex/internal/parser"

	log "github.com/sirupsen/logrus"
)

func main() {
	// we don't process error here, just parse flags
	_ = cmd.RootCmd.ParseFlags(os.Args[1:])

	if v, err := cmd.RootCmd.Flags().GetBool("verbose"); err == nil {
		if v {
			log.SetLevel(log.TraceLevel)
		}
	}

	if version, err := cmd.RootCmd.Flags().GetBool("version"); err == nil {
		if version {
			// skip register
			log.Debug("skip register user cmds, only show version")
			cmd.Execute()
			return
		}
	}

	c, err := config.ReadMakexConfig()
	if err != nil {
		log.Fatalf("failed to read makex config, err:%v", err)
	}

	makexfilename := c.Makexfile

	log.Debugf("[config] use makexfile: %s", makexfilename)

	// 2. read makexfile
	makexfile, err := parser.ReadMakexfile(makexfilename)
	if err != nil {
		log.Fatal(err)
	}

	// 3. register internal code snippet template init
	log.Debugf("[register] registering internal code commands to cobra")
	code.RegisterCmds(config.CODE_DIR_PATH, template.ExportInitCmd())

	// 4. register makexfile cmds to cobra cmds
	log.Debugf("[register] registering makexfile commands to cobra")
	makexfile.RegisterCmds(cmd.RootCmd)

	// 5. run cobra cmds executor
	log.Debugf("[execute] executing cobra")
	cmd.Execute()
}
