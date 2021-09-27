package {{.package}}

import (
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
    {{range .Columns}}{{.|firstToUpper}}   stmt.Column //{{.|firstToUpper}} is column of {{.}}
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
        {{range .Columns}}{{.|firstToUpper}}       : stmt.Column(prefix+ ".{{.}}" ) ,
        {{end}}
     }
}

func (t {{.TableName | firstToLow}}) Alias(n string) {{.TableName|firstToLow}}{
	return new{{.TableName|firstToUpper}}(t.table, n)
}

func (t {{.TableName | firstToLow}}) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

{{end}}

