package gentest

import (
	"google.golang.org/genproto/protobuf/field_mask"
)

type Inner struct {
	fieldMask field_mask.FieldMask
	Foo       string `json:"foo,omitempty"`
	Bar       int64  `json:"bar,omitempty"`
}

