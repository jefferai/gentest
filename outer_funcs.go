package gentest

import (
	"encoding/json"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/helper/strutil"
)

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
