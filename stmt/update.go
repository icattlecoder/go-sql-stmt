package stmt

import "strings"

type updateClause struct {
	tableName string
	setMap    map[Column]interface{}
	where     *baseClause
}

func Update(table Node) *updateClause {
	return &updateClause{
		tableName: table.SqlString(),
	}
}

func (c *updateClause) Set(setMap map[Column]interface{}) *updateClause {
	c.setMap = setMap
	return c
}

func (c *updateClause) Where(n ...Node) *updateClause {
	c.where = newBaseClause("WHERE", n...).setSplit(" AND ")
	return c
}

func (c *updateClause) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString("UPDATE ")
	sb.WriteString(c.tableName)

	if c.setMap != nil {
		sb.WriteString(" SET ")
		sets := make([]string, 0, len(c.setMap))
		for k, v := range c.setMap {
			set := strings.Replace(k.SqlString(), c.tableName+".", "", -1)
			set += " = "
			switch vt := v.(type) {
			case Node:
				set += vt.SqlString()
			default:
				set += "%s"
			}
			sets = append(sets, set)
		}
		sb.WriteString(strings.Join(sets, ", "))
	}
	if c.where != nil {
		sb.WriteString(" ")
		sb.WriteString(strings.Replace(c.where.SqlString(), c.tableName+".", "", -1))
	}

	return sb.String()
}

func (c *updateClause) Values() []interface{} {
	var vs []interface{}
	if c.setMap != nil {
		for _, v := range c.setMap {
			switch vt := v.(type) {
			case valueNode:
				vs = append(vs, vt.Values()...)
			default:
				vs = append(vs, v)
			}
		}
	}
	if c.where != nil {
		vs = append(vs, c.where.Values()...)
	}
	return vs
}

func (c *updateClause) writeColumns(sb *strings.Builder, columns []Column) {
	sb.WriteString("(")
	for i, col := range columns {
		sb.WriteString(strings.Replace(col.SqlString(), c.tableName+".", "", -1))
		if i < len(columns)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
}
