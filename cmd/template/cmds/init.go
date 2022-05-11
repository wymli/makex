/*
Copyright © 2022 Li Weiming <liwm29@mail2.sysu.edu.cn>

*/
package cmds

import (
	"fmt"
	"os"

	"github.com/wymli/makex/internal/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"create"},
	Short:   "init creates 'makex.yaml' in $pwd",
	Long:    `init creates 'makex.yaml' in $pwd`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. read config first
		c, err := config.ReadMakexConfig()
		if err != nil {
			log.Fatal(err)
		}

		// 2. default makexfile name is `makex.yaml`, u can change it in ~/.makex/makex_config.yaml
		makexfile := c.Makexfile

		// 3.makex.yaml init if not exists or force replace
		stat, err := os.Stat(makexfile)
		if os.IsExist(err) && stat.IsDir() {
			fmt.Printf("[init] detected existed makexfile '%s' is a dir, remove it and have a try again", makexfile)
			return
		}

		if os.IsNotExist(err) || force {
			err = config.WriteDefaultMakexfile(c, makexfile)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("makexfile detected, run `makex template init --force` to re-init")
			log.Debug("skip makexfile, exists")
		}
	},
}

var force bool

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	InitCmd.Flags().BoolVar(&force, "force", false, "force to replace makexfile in $cwd")
}
