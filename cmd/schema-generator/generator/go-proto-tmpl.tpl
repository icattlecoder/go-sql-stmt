package {{.package}}
import (
    "fmt"
)
{{range $i1, $t := .tables}}

type {{$t.TableName|firstToUpper}} struct{
    {{range .Columns}}// {{.Name}} {{.Comment}}
    {{.Name|firstToUpper}}   {{.GoType}} `{{.Tag}}`
    {{end}}
}
{{end}}

