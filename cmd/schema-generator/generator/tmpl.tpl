package {{.package}}

import (
    "strings"
    "github.com/fork-ai/go-sql-stmt/stmt"
)

var (
{{range .tables}}
    // {{.TableName | firstToUpper }} is table {{.TableName}}
    {{.TableName | firstToUpper }} =  new{{.TableName | firstToUpper}}("{{.TableName}}", "")
{{end}}
)

{{range .tables}}

type {{.TableName|firstToLow}} struct{
	table, alias string
    {{range .Columns}}{{.Name|firstToUpper}}   stmt.Column //{{.Name|firstToUpper}} is column of {{.Name}}
    {{end}}
}

func new{{.TableName|firstToUpper}}(table, alias string) {{.TableName |firstToLow }} {
     prefix := table
     if alias != "" {
         prefix = alias
     }
     return {{.TableName |firstToLow }}{
        table   : table,
        alias       :     alias,
        {{range .Columns}}{{.Name|firstToUpper}}       : stmt.Column(prefix+ ".{{.Name}}" ) ,
        {{end}}
     }
}

func (t {{.TableName | firstToLow}}) Alias(n string) {{.TableName|firstToLow}}{
	return new{{.TableName|firstToUpper}}(t.table, n)
}

func (t {{.TableName | firstToLow}}) SortableColumns(columns ...string)(n []stmt.Node){
    {{if .HasSortable}}for _, c := range columns {
		desc := strings.HasPrefix(c, "-")
		c = strings.TrimLeft(c,"-")
		{{range .Columns}}{{if .Sortable}}if c == "{{.Name}}" {
		    n = append(n, t.{{.Name|firstToUpper}}.Desc(desc))
		    continue
		}
		{{end}}{{end}}
	}{{end}}
	return
}

func (t {{.TableName | firstToLow}}) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

{{end}}

