package stmt

import "strings"

type deleteClause struct {
	tableName string
	where     *baseClause
}

func DeleteFrom(table Node) *deleteClause {
	return &deleteClause{
		tableName: table.SqlString(),
	}
}

func (c *deleteClause) Where(n ...Node) *deleteClause {
	c.where = newBaseClause("WHERE", n...).setSplit(" AND ")
	return c
}

func (c *deleteClause) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString("DELETE FROM ")
	sb.WriteString(c.tableName)

	if c.where != nil {
		sb.WriteString(" ")
		sb.WriteString(strings.Replace(c.where.SqlString(), c.tableName+".", "", -1))
	}

	return sb.String()
}

func (c *deleteClause) Values() []interface{} {
	var vs []interface{}
	if c.where != nil {
		vs = append(vs, c.where.Values()...)
	}
	return vs
}
