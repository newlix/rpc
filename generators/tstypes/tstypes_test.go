package tstypes_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/newlix/rpc/generators/tstypes"
	"github.com/newlix/rpc/schema"
)

func TestGenerate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = tstypes.Generate(&act, schema)
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_types.ts", act.Bytes())
}
