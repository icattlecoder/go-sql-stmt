package stmt

import (
	"strings"

	"github.com/keegancsmith/sqlf"
)

type updateClause struct {
	tableName string
	setMap    map[Column]interface{}
	where     *baseClause
	returning []Column
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

func (c *updateClause) Returning(columns ...Column) *updateClause {
	c.returning = columns
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
			set := k.ColumnName()
			set += " = "
			switch vt := v.(type) {
			case Column:
				set += vt.ColumnName()
			case Node:
				set += vt.SqlString()
			default:
				set += "%s"
			}
			sets = append(sets, set)
		}
		sb.WriteString(strings.Join(sets, ", "))
	}
	if c.where != nil && len(c.where.nodes) > 0 {
		sb.WriteRune(' ')
		sb.WriteString(c.where.SqlString())
	}
	if len(c.returning) > 0 {
		sb.WriteString(" RETURNING ")
		writeColumns(sb, c.returning, false)
	}

	return sb.String()
}

func (c *updateClause) Values() []interface{} {
	var vs []interface{}
	if c.setMap != nil {
		for _, v := range c.setMap {
			switch vt := v.(type) {
			case Column:
				// pass
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

func (c *updateClause) Query() *sqlf.Query {
	return sqlf.Sprintf(c.SqlString(), c.Values()...)
}
