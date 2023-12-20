package gotypes

import (
	"fmt"
	"io"
	"strings"

	"github.com/newlix/rpc/internal/format"
	"github.com/newlix/rpc/internal/schemautil"
	"github.com/newlix/rpc/schema"
)

// Generate writes the Go type implementations to w.
func Generate(w io.Writer, s *schema.Schema) error {
	out := fmt.Fprintf

	// default tags
	if s.Go.Tags == nil {
		s.Go.Tags = []string{"json"}
	}

	// types
	for _, t := range s.TypesSlice() {
		out(w, "// %s %s\n", format.GoName(t.Name), t.Description)
		out(w, "type %s struct {\n", format.GoName(t.Name))
		writeFields(w, s, t.Properties)
		out(w, "}\n\n")
	}

	// methods
	for _, m := range s.Methods {
		name := format.GoName(m.Name)

		// inputs
		if len(m.Inputs) > 0 {
			out(w, "// %sInput params.\n", name)
			out(w, "type %sInput struct {\n", name)
			writeFields(w, s, m.Inputs)
			out(w, "}\n")
		}

		// both
		if len(m.Inputs) > 0 && len(m.Outputs) > 0 {
			out(w, "\n")
		}

		// outputs
		if len(m.Outputs) > 0 {
			out(w, "// %sOutput params.\n", name)
			out(w, "type %sOutput struct {\n", name)
			writeFields(w, s, m.Outputs)
			out(w, "}\n")
		}

		out(w, "\n")
	}

	return nil
}

// writeFields to writer.
func writeFields(w io.Writer, s *schema.Schema, fields []schema.Field) {
	for i, f := range fields {
		writeField(w, s, f)
		if i < len(fields)-1 {
			fmt.Fprintf(w, "\n")
		}
	}
}

// writeField to writer.
func writeField(w io.Writer, s *schema.Schema, f schema.Field) {
	fmt.Fprintf(w, "  // %s is %s%s\n", format.GoName(f.Name), f.Description, schemautil.FormatExtra(f))
	fmt.Fprintf(w, "  %s %s %s\n", format.GoName(f.Name), goType(s, f), fieldTags(f, s.Go.Tags))
}

// goType returns a Go equivalent type for field f.
func goType(s *schema.Schema, f schema.Field) string {
	// ref
	if ref := f.Type.Ref.Value; ref != "" {
		t := schemautil.ResolveRef(s, f.Type.Ref)
		return format.GoName(t.Name)
	}

	// type
	switch f.Type.Type {
	case schema.String:
		return "string"
	case schema.Int:
		return "int"
	case schema.Bool:
		return "bool"
	case schema.Float:
		return "float64"
	case schema.Timestamp:
		return "time.Time"
	case schema.Object:
		return "map[string]interface{}"
	case schema.Array:
		return "[]" + goType(s, schema.Field{
			Type: schema.TypeObject(f.Items),
		})
	default:
		panic("unhandled type")
	}
}

// fieldTags returns tags for a field.
func fieldTags(f schema.Field, tags []string) string {
	var pairs [][]string

	for _, tag := range tags {
		pairs = append(pairs, []string{tag, f.Name})
	}

	return formatTags(pairs)
}

// formatTags returns field tags.
func formatTags(tags [][]string) string {
	var s []string
	for _, t := range tags {
		if len(t) == 2 {
			s = append(s, fmt.Sprintf("%s:%q", t[0], t[1]))
		} else {
			s = append(s, fmt.Sprintf("%s", t[0]))
		}
	}
	return fmt.Sprintf("`%s`", strings.Join(s, " "))
}
