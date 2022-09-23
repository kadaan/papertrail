package cmd

import (
	"github.com/kadaan/papertrail/lib/command"
	"github.com/kadaan/papertrail/version"
)

var (
	Root = command.NewRootCommand(
		version.Name,
		version.Name+` provides tools for interacting with `+version.Name)
)
