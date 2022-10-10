package stmt

import (
	"strings"

	"github.com/keegancsmith/sqlf"
)

type insertClause struct {
	tableName   string
	columns     []Column
	values      *valuesConstructor
	onConflict  []Column
	doUpdateSet map[Column]interface{}
}

func InsertInto(table Node, columns []Column, values Node) *insertClause {
	c := &insertClause{
		tableName: table.SqlString(),
		columns:   columns,
	}
	if values != nil {
		if vt, ok := values.(*valuesConstructor); ok {
			c.values = vt
		} else {
			c.values = Values(values)
		}
	}
	return c
}

func (c *insertClause) OnConflict(columns []Column) *insertClause {
	c.onConflict = columns
	return c
}

func (c *insertClause) DoNothing() *insertClause {
	return c
}

func (c *insertClause) DoUpdateSet(setMap map[Column]interface{}) *insertClause {
	c.doUpdateSet = setMap
	return c
}

func (c *insertClause) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString("INSERT INTO ")
	sb.WriteString(c.tableName)

	if c.columns != nil {
		sb.WriteString(" ")
		c.writeColumns(sb, c.columns)
	}
	if c.values != nil {
		sb.WriteString(" ")
		sb.WriteString(c.values.SqlString())
	}
	if c.onConflict != nil {
		sb.WriteString(" ON CONFLICT ")
		c.writeColumns(sb, c.onConflict)
		if c.doUpdateSet != nil {
			sb.WriteString(" DO UPDATE SET ")
			sets := make([]string, 0, len(c.doUpdateSet))
			for k, v := range c.doUpdateSet {
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
		} else {
			sb.WriteString(" DO NOTHING")
		}
	}

	return sb.String()
}

func (c *insertClause) Values() []interface{} {
	var vs []interface{}
	if c.values != nil {
		vs = append(vs, c.values.Values()...)
	}
	if c.onConflict != nil && c.doUpdateSet != nil {
		for _, v := range c.doUpdateSet {
			switch vt := v.(type) {
			case valueNode:
				vs = append(vs, vt.Values()...)
			default:
				vs = append(vs, v)
			}
		}
	}
	return vs
}

func (c *insertClause) writeColumns(sb *strings.Builder, columns []Column) {
	sb.WriteString("(")
	for i, col := range columns {
		sb.WriteString(strings.Replace(col.SqlString(), c.tableName+".", "", -1))
		if i < len(columns)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
}

func (c *insertClause) Query() *sqlf.Query {
	return sqlf.Sprintf(c.SqlString(), c.Values()...)
}
