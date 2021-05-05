package sqlitetypes

import (
	"fmt"
	"io"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/schema"
)

// Generate writes the Go type implementations to w, with optional validation methods.
func Generate(w io.Writer, s *schema.Schema) error {
	out := fmt.Fprintf
	out(w, "var SQLiteSchema = `\n")
	// types
	for _, t := range s.TypesSlice() {
		if !supported(s, t.Properties) {
			continue
		}
		out(w, "CREATE TABLE %s (\n", format.GoName(t.Name))
		writeFields(w, s, t.Properties)
		out(w, ");\n\n")
	}
	out(w, "`")
	return nil
}

// writeFields to writer.
func writeFields(w io.Writer, s *schema.Schema, fields []schema.Field) {
	for i, f := range fields {
		writeField(w, s, f)
		if i < len(fields)-1 {
			fmt.Fprintf(w, ",")
		}
		fmt.Fprintf(w, "\n")
	}
}

// writeField to writer.
func writeField(w io.Writer, s *schema.Schema, f schema.Field) {
	fmt.Fprintf(w, "    %s %s NOT NULL", format.GoName(f.Name), sqliteType(s, f))
	if f.Name == "id" {
		fmt.Fprintf(w, " PRIMARY KEY")
	}
}

// sqliteType returns a sqlite equivalent type for field f.
func sqliteType(s *schema.Schema, f schema.Field) string {
	// type
	switch f.Type.Type {
	case schema.String:
		return "TEXT"
	case schema.Int:
		return "INTEGER"
	case schema.Bool:
		return "INTEGER"
	case schema.Float:
		return "REAL"
	case schema.Timestamp:
		return "TEXT"
	default:
		panic("unhandled type")
	}
}

func supported(s *schema.Schema, fields []schema.Field) bool {
	for _, f := range fields {
		if f.Type.Type == schema.Array || f.Type.Type == schema.Object {
			return false
		}
	}
	return true
}
