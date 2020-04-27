package gentest

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"google.golang.org/genproto/protobuf/field_mask"
)

var _ = (json.Marshaler)((*Outer)(nil))

func init() {
	structs.DefaultTagName = "json"
}

type Outer struct {
	fieldMask field_mask.FieldMask
	Foo       *string `json:"foo,omitempty"`
	Bar       *int64  `json:"bar,omitempty"`
}

type Inner struct {
	fieldMask field_mask.FieldMask
	Foo       string `json:"foo,omitempty"`
	Bar       int64  `json:"bar,omitempty"`
}

func (o *Outer) SetDefault(key string) {
	o.fieldMask.Paths = strutil.AppendIfMissing(o.fieldMask.Paths, key)
}

func (o *Outer) UnsetDefault(key string) {
	o.fieldMask.Paths = strutil.StrListDelete(o.fieldMask.Paths, key)
}

func (o Outer) MarshalJSON() ([]byte, error) {
	m := structs.Map(o)
	if m == nil {
		m = make(map[string]interface{})
	}
	for _, k := range o.fieldMask.Paths {
		m[k] = nil
	}
	return json.Marshal(m)
}

func (i *Inner) UnmarshalJSON(in []byte) error {
	if i == nil {
		return errors.New("nil outer")
	}

	reader := strings.NewReader(string(in))

	// First decode into a map where we can check for explicit nulls and build
	// up the default map
	m := map[string]json.RawMessage{}
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(&m); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	for k, v := range m {
		if string(v) == "null" {
			i.fieldMask.Paths = strutil.AppendIfMissing(i.fieldMask.Paths, k)
		}
	}

	// Reset
	reader.Seek(0, 0)
	decoder = json.NewDecoder(reader)
	decoder.UseNumber()
	// Type aliasing is a trick here to avoid infinite recursion -- the type
	// inherits members but not methods
	type tmpInner Inner
	var ti tmpInner
	if err := decoder.Decode(&ti); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}
	// We're going to assign values to what we just decoded but we want to use
	// the fieldMask we created above, so we save it off before copying and
	// restore after
	origDefaultMap := i.fieldMask
	*i = Inner(ti)
	i.fieldMask = origDefaultMap
	return nil
}
