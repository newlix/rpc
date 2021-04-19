package swifttypes

import (
	"fmt"
	"io"
	"strings"

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
		out(w, "// %s %s\n", strcase.ToCamel(t.Name), t.Description)
		out(w, "struct %s: Codable {\n", strcase.ToCamel(t.Name))
		writeFields(w, s, t.Properties)
		if validate {
			writeValidation(w, strcase.ToLowerCamel(t.Name), t.Properties)
			out(w, "\n")
		}
		out(w, "}\n\n")
	}

	// methods
	for _, m := range s.Methods {
		name := strcase.ToCamel(m.Name)

		// inputs
		if len(m.Inputs) > 0 {
			out(w, "// %sInput params.\n", name)
			out(w, "struct %sInput: Codable {\n", name)
			writeFields(w, s, m.Inputs)
			out(w, "}\n")
			if validate {
				out(w, "\n")
				writeValidation(w, name+"Input", m.Inputs)
			}
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
	name := strcase.ToLowerCamel(f.Name)
	fmt.Fprintf(w, "    // %s is %s%s\n", name, f.Description, schemautil.FormatExtra(f))
	fmt.Fprintf(w, "    var %s: %s\n", name, swiftType(s, f))
}

// swiftType returns a Go equivalent type for field f.
func swiftType(s *schema.Schema, f schema.Field) string {
	// ref
	if ref := f.Type.Ref.Value; ref != "" {
		t := schemautil.ResolveRef(s, f.Type.Ref)
		return strcase.ToCamel(t.Name)
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

// writeValidation writes a validation method implementation to w.
func writeValidation(w io.Writer, name string, fields []schema.Field) error {
	out := fmt.Fprintf
	recv := strings.ToLower(name)[0]
	out(w, "// Validate implementation.\n")
	out(w, "func (%c *%s) Validate() error {\n", recv, name)
	for _, f := range fields {
		writeFieldDefaults(w, f, recv)
	}
	out(w, "  return nil\n")
	out(w, "}\n")
	return nil
}

// writeFieldDefaults writes field defaults to w.
func writeFieldDefaults(w io.Writer, f schema.Field, recv byte) error {
	// TODO: write out a separate Default() method?
	if f.Default == nil {
		return nil
	}

	out := fmt.Fprintf
	name := strcase.ToCamel(f.Name)

	switch f.Type.Type {
	case schema.Int:
		out(w, "  if %c.%s == 0 {\n", recv, name)
		out(w, "    %c.%s = %v\n", recv, name, f.Default)
		out(w, "  }\n\n")
	case schema.String:
		out(w, "  if %c.%s == \"\" {\n", recv, name)
		out(w, "    %c.%s = %q\n", recv, name, f.Default)
		out(w, "  }\n\n")
	}

	return nil
}
