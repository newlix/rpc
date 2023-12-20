package sqlc

import (
	"fmt"
	"io"

	"github.com/newlix/rpc/internal/format"
	"github.com/newlix/rpc/schema"
)

func supported(fs []schema.Field) []schema.Field {
	out := []schema.Field{}
	for _, f := range fs {
		if f.Type.Type != schema.Object && f.Type.Type != schema.Array {
			out = append(out, f)
		}
	}
	return out
}

func withoutID(fs []schema.Field) []schema.Field {
	out := []schema.Field{}
	for _, f := range fs {
		if f.Name != "id" {
			out = append(out, f)
		}
	}
	return out
}

// Generate writes the Go type implementations to w.
func GenerateQuery(w io.Writer, s *schema.Schema) error {
	// types
	for _, t := range s.TypesSlice() {
		fields := supported(t.Properties)
		if len(fields) == 0 {
			continue
		}
		queryGet(w, s, t)
		queryCreate(w, s, t)
		queryUpdate(w, s, t)
		queryUpsert(w, s, t)
		queryDelete(w, s, t)
	}

	return nil
}

func queryGet(w io.Writer, s *schema.Schema, t schema.Type) {
	out := fmt.Fprintf
	fields := supported(t.Properties)
	out(w, "const get%s = `SELECT ", format.GoName(t.Name))
	for i, f := range fields {
		out(w, "%q", f.Name)
		if i < len(fields)-1 {
			out(w, ", ")
		}
	}
	out(w, " FROM %q WHERE id = $1 LIMIT 1;`\n", t.Name)
	out(w, "\n")
	out(w, "func (q *Queries) Get%s(ctx context.Context, id string) (%s ,error) {\n", format.GoName(t.Name), format.GoName(t.Name))
	out(w, "	row := q.db.QueryRow(ctx, get%s, id)\n", format.GoName(t.Name))
	out(w, "	var o %s\n", format.GoName(t.Name))
	out(w, "	err := row.Scan(\n")
	for _, f := range fields {
		out(w, "		&o.%s,\n", format.GoName(f.Name))
	}
	out(w, "	)\n")
	out(w, "	return o, err\n")
	out(w, "}\n")
	out(w, "\n")
}

func queryCreate(w io.Writer, s *schema.Schema, t schema.Type) {
	out := fmt.Fprintf
	fields := supported(t.Properties)
	out(w, "const create%s = `INSERT INTO %q (", format.GoName(t.Name), t.Name)
	for i, f := range fields {
		out(w, "%q", f.Name)
		if i < len(fields)-1 {
			out(w, ", ")
		}
	}
	out(w, ")\n")
	out(w, "VALUES (")
	for i, _ := range fields {
		out(w, "$%d", i+1)
		if i < len(fields)-1 {
			out(w, ", ")
		}
	}
	out(w, ");`\n")
	out(w, "\n")
	out(w, "func (q *Queries) Create%s(ctx context.Context, o %s) error {\n", format.GoName(t.Name), format.GoName(t.Name))
	out(w, "	_, err := q.db.Exec(ctx, create%s,\n", format.GoName(t.Name))
	for _, f := range fields {
		out(w, "		o.%s,\n", format.GoName(f.Name))
	}
	out(w, "	)\n")
	out(w, "	return err\n")
	out(w, "}\n")
	out(w, "\n")
}

func queryUpsert(w io.Writer, s *schema.Schema, t schema.Type) {
	out := fmt.Fprintf
	fields := supported(t.Properties)
	out(w, "const upsert%s = `INSERT INTO %q (", format.GoName(t.Name), t.Name)
	for i, f := range fields {
		out(w, "%q", f.Name)
		if i < len(fields)-1 {
			out(w, ", ")
		}
	}
	out(w, ")\n")
	out(w, "VALUES (")
	for i, _ := range fields {
		out(w, "$%d", i+1)
		if i < len(fields)-1 {
			out(w, ", ")
		}
	}
	out(w, ")\n")
	out(w, "ON CONFLICT (\"id\")\n")
	out(w, "DO UPDATE SET ")
	for i, f := range fields {
		out(w, "%q = $%d", f.Name, i+1)
		if i < len(fields)-1 {
			out(w, ", ")
		}
	}
	out(w, ";`\n")
	out(w, "\n")
	out(w, "func (q *Queries) Upsert%s(ctx context.Context, o %s) error {\n", format.GoName(t.Name), format.GoName(t.Name))
	out(w, "	_, err := q.db.Exec(ctx, upsert%s,\n", format.GoName(t.Name))
	for _, f := range fields {
		out(w, "		o.%s,\n", format.GoName(f.Name))
	}
	out(w, "	)\n")
	out(w, "	return err\n")
	out(w, "}\n")
	out(w, "\n")
}

func queryUpdate(w io.Writer, s *schema.Schema, t schema.Type) {
	out := fmt.Fprintf
	out(w, "const update%s = `UPDATE %q\n SET ", format.GoName(t.Name), t.Name)
	fields := withoutID(supported(t.Properties))
	for i, f := range fields {
		out(w, "%q = $%d", f.Name, i+2)
		if i < len(fields)-1 {
			out(w, ",")
		}
		out(w, " ")
	}
	out(w, "WHERE id = $1;`\n")
	out(w, "\n")
	out(w, "func (q *Queries) Update%s(ctx context.Context, o %s) error {\n", format.GoName(t.Name), format.GoName(t.Name))
	out(w, "	_, err := q.db.Exec(ctx, update%s,\n", format.GoName(t.Name))
	out(w, "		o.%s,\n", format.GoName("id"))
	for _, f := range fields {
		out(w, "		o.%s,\n", format.GoName(f.Name))
	}
	out(w, "	)\n")
	out(w, "	return err\n")
	out(w, "}\n")
	out(w, "\n")
}

func queryDelete(w io.Writer, s *schema.Schema, t schema.Type) {
	out := fmt.Fprintf
	out(w, "const delete%s = `DELETE FROM %q WHERE \"id\" = $1;`\n", format.GoName(t.Name), t.Name)
	out(w, "\n")
	out(w, "func (q *Queries) Delete%s(ctx context.Context, id string) error {\n", format.GoName(t.Name))
	out(w, "	_, err := q.db.Exec(ctx, delete%s, id)\n", format.GoName(t.Name))
	out(w, "	return err\n")
	out(w, "}\n")
	out(w, "\n")
}
