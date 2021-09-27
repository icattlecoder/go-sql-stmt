package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/fork-ai/go-sql-stmt/cmd/schema-generator/generator"
)

var (
	schema string
	output string
	pkg    string
)

func init() {
	flag.StringVar(&schema, "i", "schema.json", "schema json path")
	flag.StringVar(&output, "o", "schema.go", "output schema go path")
	flag.StringVar(&pkg, "p", "", "output package name, default to sub path of output")
}

func main() {
	flag.Parse()
	output, err := filepath.Abs(output)
	if err != nil {
		log.Fatalf("invalid output: %s", output)
	}
	if pkg == "" {
		dir, _ := filepath.Split(output)
		dir = strings.TrimRightFunc(dir, func(r rune) bool {
			return r == filepath.Separator
		})
		_, pkg = filepath.Split(dir)
	}
	if err := generator.GenerateFromSchema(pkg, schema, output); err != nil {
		log.Fatalf("generate schema failed: %v, pkg: %s, schema: %s, output: %s", err, pkg, schema, output)
	}
}
