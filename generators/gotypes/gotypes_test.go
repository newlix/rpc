package gotypes_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/newlix/rpc/generators/gotypes"
	"github.com/newlix/rpc/schema"
)

func TestGenerate_noValidate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = gotypes.Generate(&act, schema)
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_types_no_validate.go", act.Bytes())
}
