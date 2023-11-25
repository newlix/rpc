package kotlintypes

import (
	"fmt"
	"io"

	"github.com/newlix/rpc/internal/schemautil"
	"github.com/newlix/rpc/schema"
	"github.com/iancoleman/strcase"
)

// Generate writes the Go type implementations to w, with optional validation methods.
func Generate(w io.Writer, s *schema.Schema, validate bool) error {
	out := fmt.Fprintf

	out(w, "import kotlinx.serialization.SerialName\n")
	out(w, "import kotlinx.serialization.Serializable\n")
	out(w, "\n")

	// types
	for _, t := range s.TypesSlice() {
		out(w, "/**\n * %s %s\n", strcase.ToCamel(t.Name), t.Description)
		writeFieldsDoc(w, s, t.Properties)
		out(w, " */\n")
		out(w, "@Serializable\n")
		out(w, "data class %s(\n", strcase.ToCamel(t.Name))
		writeFields(w, s, t.Properties)
		out(w, "\n)\n\n")
	}

	// methods
	for _, m := range s.Methods {

		// inputs
		if len(m.Inputs) > 0 {
			out(w, "/**\n * %s input params.\n", strcase.ToLowerCamel(m.Name))
			writeFieldsDoc(w, s, m.Inputs)
			out(w, " */\n")
			out(w, "@Serializable\n")
			out(w, "data class %sInput(\n", strcase.ToCamel(m.Name))
			writeFields(w, s, m.Inputs)
			out(w, "\n)\n\n")
		}

		// outputs
		if len(m.Outputs) > 0 {
			out(w, "/**\n * %s output params.\n", strcase.ToLowerCamel(m.Name))
			writeFieldsDoc(w, s, m.Outputs)
			out(w, " */\n")
			out(w, "@Serializable\n")
			out(w, "data class %sOutput(\n", strcase.ToCamel(m.Name))
			writeFields(w, s, m.Outputs)
			out(w, "\n)\n\n")
		}

	}

	return nil
}

// writeFields to writer.
func writeFieldsDoc(w io.Writer, s *schema.Schema, fields []schema.Field) {
	for _, f := range fields {
		name := strcase.ToLowerCamel(f.Name)
		fmt.Fprintf(w, " * @property %s is %s%s\n", name, f.Description, schemautil.FormatExtra(f))
	}
}

// writeFields to writer.
func writeFields(w io.Writer, s *schema.Schema, fields []schema.Field) {
	for i, f := range fields {
		writeField(w, s, f)
		if i < len(fields)-1 {
			fmt.Fprintf(w, ",\n")
		}
	}
}

// writeField to writer.
func writeField(w io.Writer, s *schema.Schema, f schema.Field) {
	t := "var"
	if f.ReadOnly {
		t = "val"
	}
	fmt.Fprintf(w, "    @SerialName(\"%s\") %s %s: %s = %s", f.Name, t, strcase.ToLowerCamel(f.Name), kotlinType(s, f), defaultValue(s, f))
}

// kotlinType returns a Kotlin equivalent type for field f.
func kotlinType(s *schema.Schema, f schema.Field) string {
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
		return "Boolean"
	case schema.Float:
		return "Double"
	case schema.Timestamp:
		return "String"
	case schema.Object:
		return kotlinType(s, schema.Field{
			Type: schema.TypeObject(f.Items),
		})
	case schema.Array:
		return "Array<" + kotlinType(s, schema.Field{
			Type: schema.TypeObject(f.Items),
		}) + ">"
	default:
		panic("unhandled type")
	}
}

func defaultValue(s *schema.Schema, f schema.Field) string {
	if ref := f.Type.Ref.Value; ref != "" {
		t := schemautil.ResolveRef(s, f.Type.Ref)
		return strcase.ToCamel(t.Name) + "()"
	}

	// type
	switch f.Type.Type {
	case schema.String:
		return "\"\""
	case schema.Int:
		return "0"
	case schema.Bool:
		return "false"
	case schema.Float:
		return "0.0"
	case schema.Timestamp:
		return "\"\""
	case schema.Object:
		return kotlinType(s, schema.Field{
			Type: schema.TypeObject(f.Items),
		}) + "()"
	case schema.Array:
		return "arrayOf()"
	default:
		panic("unhandled type")
	}
}
