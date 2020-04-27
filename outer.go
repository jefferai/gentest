package gentest

import (
	"encoding/json"

	"google.golang.org/genproto/protobuf/field_mask"
)

var _ = (json.Marshaler)((*Outer)(nil))

type Outer struct {
	fieldMask field_mask.FieldMask
	Foo       *string `json:"foo,omitempty"`
	Bar       *int64  `json:"bar,omitempty"`
}
