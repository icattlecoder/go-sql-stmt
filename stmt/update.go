package stmt

import (
	"sort"
	"strings"

	"github.com/keegancsmith/sqlf"
)

type updateClause struct {
	tableName string
	setMap    map[Column]interface{}
	setCols   []Column
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
	c.setCols = make([]Column, 0, len(setMap))
	for k := range setMap {
		c.setCols = append(c.setCols, k)
	}
	sort.Slice(c.setCols, func(i, j int) bool {
		return c.setCols[i].ColumnName() < c.setCols[j].ColumnName()
	})
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
		for _, column := range c.setCols {
			v := c.setMap[column]
			set := column.ColumnName()
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
		for _, column := range c.setCols {
			v := c.setMap[column]
			switch vt := v.(type) {
			case Column:
				// do nothing
			case valueNode:
				vs = append(vs, vt.Values()...)
			default:
				vs = append(vs, vt)
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
