package swifttypes_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/newlix/rpc/generators/swifttypes"
	"github.com/newlix/rpc/schema"
)

func TestGenerate_noValidate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = swifttypes.Generate(&act, schema, false)
	assert.NoError(t, err, "generating")
	// os.WriteFile("/tmp/dat1", act.Bytes(), 0644)
	fixture.Assert(t, "todo_types_no_validate.swift", act.Bytes())
}
