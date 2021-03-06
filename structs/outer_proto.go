// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-04-28 22:28:25.4674158 -0400 EDT m=+0.001894701
package structs

import (
	"encoding/json"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/helper/strutil"
)

type OuterProto struct {
	defaultFields []string

	Uri          *string `json:"uri,omitempty"`
	FriendlyName *string `json:"friendly_name,omitempty"`
	// The time this host was created.
	CreatedTime time.Time `json:"created_time,omitempty"`
	// The time this host was last updated.
	UpdatedTime time.Time `json:"updated_time,omitempty"`
	// Marks the host as disabled.  Default is false.
	Disabled *bool `json:"disabled,omitempty"`
	// This field is required.
	Address *string `json:"address,omitempty"`
}

func (s *OuterProto) SetDefault(key string) {
	s.defaultFields = strutil.AppendIfMissing(s.defaultFields, key)
}

func (s *OuterProto) UnsetDefault(key string) {
	s.defaultFields = strutil.StrListDelete(s.defaultFields, key)
}

func (s OuterProto) MarshalJSON() ([]byte, error) {
	m := structs.Map(s)
	if m == nil {
		m = make(map[string]interface{})
	}
	for _, k := range s.defaultFields {
		m[k] = nil
	}
	return json.Marshal(m)
}
