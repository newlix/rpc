package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/apex/rpc/generators/swifttypes"
	"github.com/apex/rpc/schema"
)

func main() {
	path := flag.String("schema", "schema.json", "Path to the schema file")
	flag.Parse()

	s, err := schema.Load(*path)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = generate(os.Stdout, s)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// generate implementation.
func generate(w io.Writer, s *schema.Schema) error {
	err := swifttypes.Generate(w, s, false)
	if err != nil {
		return fmt.Errorf("generating types: %w", err)
	}

	return nil
}
