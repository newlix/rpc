package sqlitetypes_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/apex/rpc/generators/sqlitetypes"
	"github.com/apex/rpc/schema"
)

func TestGenerate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = sqlitetypes.Generate(&act, schema)
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_sqlite_schema.go", act.Bytes())
}
