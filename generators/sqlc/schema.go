package sqlc

import (
	"fmt"
	"io"

	"github.com/newlix/rpc/internal/format"
	"github.com/newlix/rpc/internal/schemautil"
	"github.com/newlix/rpc/schema"
)

// Generate writes the Go type implementations to w.
func GenerateSchema(w io.Writer, s *schema.Schema) error {
	out := fmt.Fprintf

	// default tags
	if s.Go.Tags == nil {
		s.Go.Tags = []string{"json"}
	}

	// types
	for _, t := range s.TypesSlice() {
		fields := supported(t.Properties)
		if len(fields) == 0 {
			continue
		}
		out(w, "-- %s %s\n", t.Name, t.Description)
		out(w, "CREATE TABLE IF NOT EXISTS %q (\n", t.Name)
		out(w, "  id   text PRIMARY KEY\n")
		out(w, ");\n")
		writeFields(w, s, withoutID(fields), t.Name)
		out(w, "\n")
	}

	return nil
}

// writeFields to writer.
func writeFields(w io.Writer, s *schema.Schema, fields []schema.Field, table string) {
	out := fmt.Fprintf
	for _, f := range fields {
		out(w, "ALTER TABLE %q ADD COLUMN IF NOT EXISTS %q %s;\n", table, f.Name, sqlType(s, f))
		out(w, "ALTER TABLE %q ALTER COLUMN %q SET NOT NULL;\n", table, f.Name)
	}
	out(w, "\n")
}

// sqlType returns a SQL equivalent type for field f.
func sqlType(s *schema.Schema, f schema.Field) string {
	// ref
	if ref := f.Type.Ref.Value; ref != "" {
		t := schemautil.ResolveRef(s, f.Type.Ref)
		return format.GoName(t.Name)
	}

	// type
	switch f.Type.Type {
	case schema.String:
		return "TEXT"
	case schema.Int:
		return "BIGINT"
	case schema.Bool:
		return "BOOLEAN"
	case schema.Float:
		return "FLOAT"
	case schema.Timestamp:
		return "DATE"
	default:
		panic("unhandled type")
	}
}
