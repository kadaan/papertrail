package config

import "time"

type SearchConfig struct {
	SystemID       uint
	GroupID        uint
	Start          time.Time
	End            time.Time
	Fields         []FieldType
	FieldSeparator string
	Limit          uint
}
