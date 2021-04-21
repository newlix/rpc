package swifttypes

import (
	"fmt"
	"io"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/internal/schemautil"
	"github.com/apex/rpc/schema"
	"github.com/iancoleman/strcase"
)

// Generate writes the Go type implementations to w, with optional validation methods.
func Generate(w io.Writer, s *schema.Schema, validate bool) error {
	out := fmt.Fprintf
	out(w, "import Foundation\n")
	out(w, "\n")

	// types
	for _, t := range s.TypesSlice() {
		out(w, "// %s %s\n", format.GoName(t.Name), t.Description)
		out(w, "struct %s: Codable {\n", format.GoName(t.Name))
		writeFields(w, s, t.Properties)
		out(w, "\n")
		writeCodingKeys(w, s, t.Properties)

		out(w, "}\n\n")
	}

	// methods
	for _, m := range s.Methods {
		name := format.GoName(m.Name)

		// inputs
		if len(m.Inputs) > 0 {
			out(w, "// %sInput params.\n", name)
			out(w, "struct %sInput: Codable {\n", name)
			writeFields(w, s, m.Inputs)
			out(w, "\n")
			writeCodingKeys(w, s, m.Inputs)
			out(w, "}\n")
		}

		// both
		if len(m.Inputs) > 0 && len(m.Outputs) > 0 {
			out(w, "\n")
		}

		// outputs
		if len(m.Outputs) > 0 {
			out(w, "// %sOutput params.\n", name)
			out(w, "struct %sOutput: Codable {\n", name)
			writeFields(w, s, m.Outputs)
			out(w, "\n")
			writeCodingKeys(w, s, m.Outputs)
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
	name := strcase.ToLowerCamel(format.GoName(f.Name))
	fmt.Fprintf(w, "    // %s is %s%s\n", name, f.Description, schemautil.FormatExtra(f))
	fmt.Fprintf(w, "    var %s: %s\n", name, swiftType(s, f))
}

// writeCodingKeys to writer
func writeCodingKeys(w io.Writer, s *schema.Schema, fields []schema.Field) {
	fmt.Fprintf(w, "    enum CodingKeys: String, CodingKey {\n")
	for _, f := range fields {
		writeCodingKey(w, s, f)
	}
	fmt.Fprintf(w, "    }\n")
}

// writeCodingKeys to writer
func writeCodingKey(w io.Writer, s *schema.Schema, f schema.Field) {
	fmt.Fprintf(w, "        case %s = \"%s\"\n", strcase.ToLowerCamel(format.GoName(f.Name)), f.Name)
}

// swiftType returns a Go equivalent type for field f.
func swiftType(s *schema.Schema, f schema.Field) string {
	// ref
	if ref := f.Type.Ref.Value; ref != "" {
		t := schemautil.ResolveRef(s, f.Type.Ref)
		return format.GoName(t.Name)
	}

	// type
	switch f.Type.Type {
	case schema.String:
		return "String"
	case schema.Int:
		return "Int"
	case schema.Bool:
		return "Bool"
	case schema.Float:
		return "Double"
	case schema.Timestamp:
		return "Date"
	case schema.Object:
		return swiftType(s, schema.Field{
			Type: schema.TypeObject(f.Items),
		})
	case schema.Array:
		return "[" + swiftType(s, schema.Field{
			Type: schema.TypeObject(f.Items),
		}) + "]"
	default:
		panic("unhandled type")
	}
}

//TODO validation
