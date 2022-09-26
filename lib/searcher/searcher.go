package searcher

import (
	"fmt"
	"github.com/kadaan/papertrail/config"
	"github.com/kadaan/papertrail/lib/command"
	"github.com/kadaan/papertrail/lib/errors"
	"github.com/kadaan/papertrail/lib/papertrail"
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

	events := searchResp.Events
	if len(events) > 0 && cfg.Limit <= uint(len(events)) {
		events = events[len(events)-int(cfg.Limit):]
	}

	var b strings.Builder
	for _, e := range events {
		for _, field := range cfg.Fields {
			switch field {
			case config.ReceivedAt:
				writeField(cfg, &b, &e.ReceivedAt)
				break
			case config.SourceName:
				writeField(cfg, &b, &e.SourceName)
				break
			case config.SourceIP:
				writeField(cfg, &b, &e.SourceIP)
				break
			case config.Facility:
				writeField(cfg, &b, &e.Facility)
				break
			case config.Program:
				writeField(cfg, &b, e.Program)
				break
			case config.Message:
				writeField(cfg, &b, &e.Message)
				break
			}
		}
		if b.Len() > 0 {
			fmt.Println(b.String())
		}
		b.Reset()
	}
	return nil
}

func writeField[T any](cfg *config.SearchConfig, stringBuilder *strings.Builder, value *T) {
	if stringBuilder.Len() > 0 {
		fieldSeparator := cfg.FieldSeparator
		if len(fieldSeparator) == 0 {
			fieldSeparator = "\x00"
		}
		stringBuilder.WriteString(fieldSeparator)
	}
	if value != nil {
		stringBuilder.WriteString(fmt.Sprintf("%s", *value))
	}
}
