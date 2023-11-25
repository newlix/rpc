package goserver_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/newlix/rpc/generators/goserver"
	"github.com/newlix/rpc/schema"
)

func TestGenerate_noTypes(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = goserver.Generate(&act, schema, "")
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_server_no_types.go", act.Bytes())
}
func TestGenerate_types(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = goserver.Generate(&act, schema, "api")
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_server_types.go", act.Bytes())
}
