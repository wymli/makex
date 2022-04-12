/*
Copyright © 2022 Li Weiming <liwm29@mail2.sysu.edu.cn>

*/
package cmds

import (
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
		c, err := config.ReadMakexConfig()
		if err != nil {
			log.Fatal(err)
		}

		if err = config.MoveShells(); err != nil {
			log.Fatal(err)
		}

		makexfile := viper.GetString("makexfile")

		err = config.WriteMakexfile(c, makexfile)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
