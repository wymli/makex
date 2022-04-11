package template

import (
	"github.com/spf13/cobra"
	"github.com/wymli/makex/cmd/template/cmds"
)

var TemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "use 'makex template init'",
	Long:  `use 'makex template init'`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	TemplateCmd.AddCommand(cmds.InitCmd)
}
