package shell

import (
	"embed"
)

//go:embed *.sh
var ShellFS embed.FS
