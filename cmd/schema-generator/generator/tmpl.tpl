package {{.package}}

import (
    "strings"

    "github.com/fork-ai/go-sql-stmt/stmt"
)

var (
{{range $i1, $t := .tables}}
    // {{$t.TableName | firstToUpper }} is table {{$t.TableName}}
    {{$t.TableName | firstToUpper }} =  New{{$t.TableName | firstToUpper}}("{{$t.TableName}}", "")
{{end}}
)

{{range $i1, $t := .tables}}

type {{$t.TableName|firstToLow}} struct{
	table, alias string
    {{range .Columns}}{{.Name|firstToUpper}}   stmt.Column //{{.Name|firstToUpper}} is column of {{.Name}}
    {{end}}
}

func New{{$t.TableName|firstToUpper}}(table, alias string) {{$t.TableName |firstToLow }} {
     prefix := table
     if alias != "" {
         prefix = alias
     }
     return {{$t.TableName |firstToLow }}{
        table   : table,
        alias       :     alias,
        {{range .Columns}}{{.Name|firstToUpper}}       : stmt.Column(prefix+ ".{{.Name}}" ) ,
        {{end}}
     }
}

func (t {{$t.TableName | firstToLow}}) All() []stmt.Node{
    return []stmt.Node{ {{range .Columns}}t.{{.Name|firstToUpper}},{{end}}}
}

func (t {{$t.TableName | firstToLow}}) Alias(n string) {{$t.TableName|firstToLow}}{
	return New{{$t.TableName|firstToUpper}}(t.table, n)
}

func (t {{$t.TableName | firstToLow}}) SortableColumns(columns ...string)(n []stmt.Node){
    {{if .HasSortable}}for _, c := range columns {
		desc := strings.HasPrefix(c, "-")
		c = strings.TrimLeft(c,"-")
		{{range .Columns}}{{if .Sortable}}if c == "{{.Name}}" {
		    n = append(n, t.{{.Name|firstToUpper}}.Desc(desc))
		    continue
		}
		{{end}}{{end}}
	}{{end}}
	return n
}

func (t {{$t.TableName | firstToLow}}) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

{{if $t.Update}}
type Update{{$t.TableName|firstToUpper}}Option struct {
	setMap map[stmt.Column]interface{}
}

func (o Update{{$t.TableName|firstToUpper}}Option) GetSetMap() map[stmt.Column]interface{}{
    return o.setMap
}

type Update{{$t.TableName|firstToUpper}}OptionFunc func(opt *Update{{$t.TableName|firstToUpper}}Option)

func WithUpdate{{$t.TableName|firstToUpper}}OptionFunc(opts ...Update{{$t.TableName|firstToUpper}}OptionFunc) *Update{{$t.TableName|firstToUpper}}Option {
	option := Update{{$t.TableName|firstToUpper}}Option{setMap: map[stmt.Column]interface{}{ {{if $t.HasUpdateAt}}"updated_at": stmt.Now(){{end}} }}
	for _, f := range opts {
		f(&option)
	}
	return &option
}

{{range $i, $c := .Columns}}{{if $c.Update}}
func WithUpdate{{$t.TableName|firstToUpper}}{{$c.Name|firstToUpper}}({{$c.Name|firstToLow}} {{$c.GoType}}) Update{{$t.TableName|firstToUpper}}OptionFunc {
	return func(opt *Update{{$t.TableName|firstToUpper}}Option) {
		opt.setMap[{{$t.TableName|firstToUpper}}.{{$c.Name|firstToUpper}}] = {{$c.Name|firstToLow}}
	}
}
{{end}}{{end}}
{{end}}
{{end}}

