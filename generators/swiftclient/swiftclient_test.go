package swiftclient_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/newlix/rpc/generators/swiftclient"
	"github.com/newlix/rpc/schema"
)

func TestGenerate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = swiftclient.Generate(&act, schema, "Client")
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_client.swift", act.Bytes())
}
