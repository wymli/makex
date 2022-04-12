/*
Copyright © 2022 Li Weiming <liwm29@mail2.sysu.edu.cn>

*/
package cmds

import (
	"os"

	"github.com/wymli/makex/internal/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"create"},
	Short:   "init  creates the '.makex' config dir in $HOME, and creates 'makex.yaml' in $pwd",
	Long:    `init  creates the '.makex' config dir in $HOME, and creates 'makex.yaml' in $pwd`,
	Run: func(cmd *cobra.Command, args []string) {
		// if err := config.InitMakexConfig(); err != nil {
		// 	log.Fatal(err)
		// }

		c, err := config.ReadMakexConfig()
		if err != nil {
			log.Fatal(err)
		}

		makexfile := viper.GetString(config.MAKEXFILE_KEY)

		// makex.yaml init if not exists or force replace
		_, err = os.Stat(makexfile)
		if os.IsNotExist(err) || force {
			err = config.WriteMakexfile(c, makexfile)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Info("makexfile detected, run `makex template init --force` to re-init")
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
