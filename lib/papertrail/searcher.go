package papertrail

import (
	"github.com/kadaan/papertrail/lib/errors"
)

type Searcher interface {
	Search(searchOptions SearchOptions, limit uint) ([]*Event, error)
}

func NewSearcher() Searcher {
	return &searcher{}
}

type searcher struct {
}

func (s searcher) Search(searchOptions SearchOptions, limit uint) ([]*Event, error) {
	token, err := ReadToken()
	if err == ErrNoTokenFound {
		return nil, errors.New("No Papertrail API token found.\n\npapertrail requires a " +
			"valid Papertrail API token (which you can obtain from https://papertrailapp.com/user/edit)\nto be set " +
			"in the PAPERTRAIL_API_TOKEN environment variable or in ~/.papertrail.yml (in the format `token: MYTOKEN`).")
	} else if err != nil {
		return nil, err
	}

	client := NewClient((&TokenTransport{Token: token}).Client())

	searchResp, _, err := client.Search(searchOptions)
	if err != nil {
		return nil, err
	}

	events := searchResp.Events
	if len(events) > 0 && limit <= uint(len(events)) {
		events = events[len(events)-int(limit):]
	}

	return events, nil
}
