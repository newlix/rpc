package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/newlix/rpc/generators/goclient"
	"github.com/newlix/rpc/generators/gotypes"
	"github.com/newlix/rpc/schema"
)

func main() {
	path := flag.String("schema", "schema.json", "Path to the schema file")
	pkg := flag.String("package", "client", "Name of the package")
	flag.Parse()

	s, err := schema.Load(*path)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = generate(os.Stdout, s, *pkg)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// generate implementation.
func generate(w io.Writer, s *schema.Schema, pkg string) error {
	out := fmt.Fprintf

	// force tags to be json only
	s.Go.Tags = []string{"json"}

	out(w, "// Do not edit, this file was generated by github.com/newlix/rpc.\n\n")
	out(w, "package %s\n\n", pkg)

	out(w, "import (\n")
	out(w, "  \"bytes\"\n")
	out(w, "  \"encoding/json\"\n")
	out(w, "  \"fmt\"\n")
	out(w, "  \"io\"\n")
	out(w, "  \"net/http\"\n")
	out(w, "  \"time\"\n")
	out(w, ")\n\n")

	err := gotypes.Generate(w, s)
	if err != nil {
		return fmt.Errorf("generating types: %w", err)
	}

	err = goclient.Generate(w, s)
	if err != nil {
		return fmt.Errorf("generating client: %w", err)
	}

	return nil
}
