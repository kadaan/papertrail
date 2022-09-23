package searcher

import (
	"fmt"
	"github.com/kadaan/papertrail/config"
	"github.com/kadaan/papertrail/lib/command"
	"github.com/kadaan/papertrail/lib/errors"
	"github.com/kadaan/papertrail/lib/papertrail"
	"k8s.io/klog/v2"
	"strings"
)

func NewSearcher() command.Task[config.SearchConfig] {
	return &searcher{}
}

type searcher struct {
}

func (v *searcher) Run(cfg *config.SearchConfig, args []string) error {
	if len(args) == 0 {
		return errors.New("No search criteria provided")
	}

	token, err := papertrail.ReadToken()
	if err == papertrail.ErrNoTokenFound {
		return errors.NewCommandError("No Papertrail API token found.\n\npapertrail requires a " +
			"valid Papertrail API token (which you can obtain from https://papertrailapp.com/user/edit)\nto be set " +
			"in the PAPERTRAIL_API_TOKEN environment variable or in ~/.papertrail.yml (in the format `token: MYTOKEN`).")
	} else if err != nil {
		return err
	}

	client := papertrail.NewClient((&papertrail.TokenTransport{Token: token}).Client())

	searchOptions := papertrail.SearchOptions{
		MinTime: cfg.Start,
		MaxTime: cfg.End,
		Query:   strings.Join(args, " "),
	}

	if cfg.GroupID > 0 {
		searchOptions.GroupID = fmt.Sprintf("%d", cfg.GroupID)
	}
	if cfg.SystemID > 0 {
		searchOptions.SystemID = fmt.Sprintf("%d", cfg.SystemID)
	}

	searchResp, _, err := client.Search(searchOptions)
	if err != nil {
		return err
	}

	for _, e := range searchResp.Events {
		var prog string
		if e.Program != nil {
			prog = *e.Program
		}
		klog.V(0).Infof("%s %s %s %s: %s\n", e.ReceivedAt, e.SourceName, e.Facility, prog, e.Message)
	}
	return nil
}
