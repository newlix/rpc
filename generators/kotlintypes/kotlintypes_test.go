package kotlintypes_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/newlix/rpc/generators/kotlintypes"
	"github.com/newlix/rpc/schema"
)

func TestGenerate_noValidate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = kotlintypes.Generate(&act, schema, false)
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_types_no_validate.kt", act.Bytes())
}
