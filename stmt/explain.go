package stmt

import (
	"fmt"
	"strings"
)

var validFormats = map[string]struct{}{
	"XML":  {},
	"JSON": {},
	"YAML": {},
}

type explain struct {
	format  string
	Analyze bool
	sel     *selectClause
	update  *updateClause
	insert  *insertClause
	del     *deleteClause
}

func (e *explain) JsonFormat() *explain {
	e.format = "JSON"
	return e
}

func (e *explain) XmlFormat() *explain {
	e.format = "XML"
	return e
}

func (e *explain) YamlFormat() *explain {
	e.format = "YAML"
	return e

}

func (e *explain) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString("EXPLAIN")

	var opts []string
	if _, ok := validFormats[e.format]; ok {
		opts = append(opts, fmt.Sprintf("FORMAT %s", e.format))
	}

	if e.Analyze {
		opts = append(opts, "ANALYZE")
	}

	if len(opts) > 0 {
		sb.WriteString(" (")
		sb.WriteString(strings.Join(opts, ","))
		sb.WriteString(")")
	}

	sb.WriteString(" ")
	if e.sel != nil {
		sb.WriteString(e.sel.SqlString())
	} else if e.update != nil {
		sb.WriteString(e.update.SqlString())
	} else if e.insert != nil {
		sb.WriteString(e.insert.SqlString())
	} else if e.del != nil {
		sb.WriteString(e.del.SqlString())
	}
	return sb.String()
}

func (e *explain) Values() []interface{} {
	var values []interface{}
	if e.sel != nil {
		values = e.sel.Values()
	} else if e.update != nil {
		values = e.update.Values()
	} else if e.insert != nil {
		values = e.insert.Values()
	} else if e.del != nil {
		values = e.del.Values()
	}
	return values
}
