package stmt

import (
	"github.com/keegancsmith/sqlf"
)

type deleteClause struct {
	tableName string
	using     *baseClause
	where     *baseClause
}

func DeleteFrom(table Node) *deleteClause {
	return &deleteClause{
		tableName: table.SqlString(),
	}
}

func (c *deleteClause) Using(n ...Node) *deleteClause {
	c.using = newBaseClause("USING", n...).setSplit(", ")
	return c
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

	if c.using != nil && len(c.using.nodes) > 0 {
		sb.WriteRune(' ')
		sb.WriteString(c.using.SqlString())
	}
	if c.where != nil && len(c.where.nodes) > 0 {
		sb.WriteRune(' ')
		sb.WriteString(c.where.SqlString())
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

func (c *deleteClause) Query() *sqlf.Query {
	return sqlf.Sprintf(c.SqlString(), c.Values()...)
}
