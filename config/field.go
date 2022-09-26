package config

import "github.com/thediveo/enumflag/v2"

type FieldType enumflag.Flag

const (
	ReceivedAt FieldType = iota
	SourceName
	SourceIP
	Facility
	Program
	Message
)

var FieldTypeIds = map[FieldType][]string{
	ReceivedAt: {"received_at"},
	SourceName: {"source_name"},
	SourceIP:   {"source_ip"},
	Facility:   {"facility"},
	Program:    {"program"},
	Message:    {"message"},
}
