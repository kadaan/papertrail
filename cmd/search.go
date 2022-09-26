package cmd

import (
	"github.com/kadaan/papertrail/config"
	"github.com/kadaan/papertrail/lib/command"
	"github.com/kadaan/papertrail/lib/searcher"
	"github.com/spf13/cobra"
)

func init() {
	command.NewCommand(
		Root,
		"search",
		"search <QUERY>",
		"Search papertrail",
		"Search papertrail for specific logs for the specified targets.",
		new(config.SearchConfig),
		searcher.NewSearcher()).Configure(func(cb command.CommandBuilder, fb config.FlagBuilder, cfg *config.SearchConfig) {
		cb.Args(cobra.MinimumNArgs(1))
		fb.TimeRange(&cfg.Start, &cfg.End, "time range to search")
		fb.GroupID(&cfg.GroupID, "group id to limit search within")
		fb.SystemID(&cfg.SystemID, "system id to search within")
		fb.Fields(&cfg.Fields, "which fields to include in the result")
		fb.FieldSeparator(&cfg.FieldSeparator, "string used to separate field values")
		fb.Limit(&cfg.Limit, "maximum number of results to return")
	})
}
