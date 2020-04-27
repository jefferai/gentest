package gentest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Marshal(t *testing.T) {
	var o Outer
	o.Foo = String("zip")

	b, err := json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, `{"foo":"zip"}`, string(b))

	o.Foo = nil
	o.Bar = Int(45)
	b, err = json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, `{"bar":45}`, string(b))

	o.SetDefault("zap")
	b, err = json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, `{"bar":45,"zap":null}`, string(b))

	o.UnsetDefault("zap")
	b, err = json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, `{"bar":45}`, string(b))
}

func Test_Unmarshal(t *testing.T) {
	var exp Inner
	var test Inner

	{
		input := `{"foo":"zip"}`
		exp.Foo = "zip"

		err := json.Unmarshal([]byte(input), &test)
		assert.NoError(t, err)
		assert.Equal(t, exp, test)
	}

	{
		input := `{"bar":45,"zap":null}`
		exp.Foo = ""
		exp.Bar = 45
		exp.fieldMask.Paths = []string{"zap"}
		test = Inner{}
		err := json.Unmarshal([]byte(input), &test)
		assert.NoError(t, err)
		assert.Equal(t, exp, test)
	}
}
