package sqlc

import (
	"fmt"
	"io"

	"github.com/newlix/rpc/internal/format"
	"github.com/newlix/rpc/schema"
)

func GenerateMigrate(w io.Writer, s *schema.Schema) error {

	// generateExport(w, s.TypesSlice())
	generateLoad(w, s.TypesSlice())

	return nil
}

func generateExport(w io.Writer, ts []schema.Type) {
	out := fmt.Fprintf
	out(w, "func export(db *gorm.DB) {\n")
	for _, t := range ts {
		if t.Name == "index" || t.Name == "index2" {
			continue
		}
		out(w, "	%ss := []model.%s{}\n", t.Name, format.GoName(t.Name))
		out(w, "	if r := db.Find(&%ss); r.Error != nil {\n", t.Name)
		out(w, "		log.Fatal(r.Error)\n")
		out(w, "	}\n")
		out(w, "	%sf, err := os.Create(\"%s.json\")\n", t.Name, t.Name)
		out(w, "	if err != nil {\n")
		out(w, "		log.Fatal(err)\n")
		out(w, "	}\n")
		out(w, "	defer %sf.Close()\n", t.Name)
		out(w, "	if err := json.NewEncoder(%sf).Encode(%ss); err != nil {\n", t.Name, t.Name)
		out(w, "		log.Fatal(err)\n")
		out(w, "	}\n")
		out(w, "\n")
	}
	out(w, "}\n")
}

func generateLoad(w io.Writer, ts []schema.Type) {
	out := fmt.Fprintf
	out(w, "func load(ctx context.Context, db *gorm.DB) {\n")
	for _, t := range ts {
		if t.Name == "index" || t.Name == "index2" {
			continue
		}
		out(w, "%sf, err := os.Open(\"%s.json\")\n", t.Name, t.Name)
		out(w, "if err != nil {\n")
		out(w, "	log.Fatal(err)\n")
		out(w, "}\n")
		out(w, "%ss := []model.%s{}\n", t.Name, format.GoName(t.Name))
		out(w, "if err := json.NewDecoder(%sf).Decode(&%ss); err != nil {\n", t.Name, t.Name)
		out(w, "	log.Fatal(err)\n")
		out(w, "}\n")
		out(w, "if r := db.Create(&%ss); r.Error != nil {\n", t.Name)
		out(w, "	log.Fatal(r.Error)\n")
		out(w, "}\n")
	}
	out(w, "}\n")
}
