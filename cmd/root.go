/*
Copyright © 2022 Li Weiming <liwm29@mail2.sysu.edu.cn>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/wymli/makex/cmd/template"
	"github.com/wymli/makex/internal/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootFlags describes a struct that holds flags that can be set on root level of the command
type RootFlags struct {
	version bool
}

type GlobalFlags struct {
	makexfile string
	verbose   bool
}

var (
	rootFlags = RootFlags{}
	globFlags = GlobalFlags{}
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "makex",
	Short: "makex is a cmd-line tool like make and task",
	Long:  `makex is a cmd-line tool like make and task`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if rootFlags.version {
			printVersion()
		} else {
			cmd.Usage()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initViper()

	// initConfig func will be called after flag parse and before cmd execute
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVarP(&globFlags.makexfile, "makexfile", "f", "", "makexfile name in cwd, default is makex.yaml")
	RootCmd.PersistentFlags().BoolVarP(&globFlags.verbose, "verbose", "v", false, "show debug log or not")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolVar(&rootFlags.version, "version", false, "show makex version")

	RootCmd.AddCommand(template.TemplateCmd)
}

func initConfig() {
}

func initViper() {
	// use default
	viper.SetDefault("makexfile", "makex.yaml")

	// use env
	viper.BindEnv("makexfile", "MAKEXFILE", "makexfile")

	// use config file
	viper.SetConfigName(config.CONFIG_NAME)
	viper.SetConfigType(config.CONFIG_TYPE)
	viper.AddConfigPath(config.HOME)

	if err := viper.ReadInConfig(); err != nil {
		log.Debugf("failed to read config file %s, call `makex template init` first, err: %v", config.CONFIG_PATH, err)
	}
}

func printVersion() {
	fmt.Println("version will be set using build-flags, and printed here")
}
