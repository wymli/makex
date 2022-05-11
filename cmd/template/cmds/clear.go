/*
Copyright © 2022 Li Weiming <liwm29@mail2.sysu.edu.cn>

*/
package cmds

import (
	"os"

	"github.com/wymli/makex/internal/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"remove", "delete"},
	Short:   "clear removes the '.makex' config dir in $HOME, but remains 'makex.yaml' in $pwd",
	Long:    `clear removes the '.makex' config dir in $HOME, but remains 'makex.yaml' in $pwd`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("[remove] clearing config dir '%s'", config.CONFIG_DIR_PATH)
		if err := os.RemoveAll(config.CONFIG_DIR_PATH); err != nil {
			log.Fatalf("failed to remove makex config dir '%s', err: %v", config.CONFIG_DIR_PATH, err)
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
