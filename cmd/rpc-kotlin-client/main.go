package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/apex/rpc/generators/kotlinclient"
	"github.com/apex/rpc/schema"
)

func main() {
	path := flag.String("schema", "schema.json", "Path to the schema file")
	pkg := flag.String("package", "com.example", "Name of the package")
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
	out(w, "// Do not edit, this file was generated by github.com/apex/rpc.\n\n")
	out(w, "package %s\n\n", pkg)
	err := kotlinclient.Generate(w, s)
	if err != nil {
		return fmt.Errorf("generating client: %w", err)
	}

	return nil
}
