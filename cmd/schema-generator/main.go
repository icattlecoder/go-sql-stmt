package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/icattlecoder/go-sql-stmt/cmd/schema-generator/generator"
)

var (
	schema string
	output string
	pkg    string
	cmd    string
)

func init() {
	flag.StringVar(&schema, "i", "schema.json", "schema json path")
	flag.StringVar(&output, "o", "schema.go", "output schema go path")
	flag.StringVar(&cmd, "c", "schema", "生成器")
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
	if cmd == "schema" {
		if err := generator.GenerateSchemaFromSchema(pkg, schema, output); err != nil {
			log.Fatalf("generate schema failed: %v, pkg: %s, schema: %s, output: %s", err, pkg, schema, output)
		}
		return
	}
	if err := generator.GenerateGoProtoFromSchema(pkg, schema, output); err != nil {
		log.Fatalf("generate schema failed: %v, pkg: %s, schema: %s, output: %s", err, pkg, schema, output)
	}
}
