package stmt

import (
	"strings"

	"github.com/keegancsmith/sqlf"
)

type insertClause struct {
	tableName      string
	columns        []Column
	values         *valuesConstructor
	onConflict     bool
	conflictTarget []Column
	onConstraint   Node
	doUpdateSet    map[Column]interface{}
	doUpdateWhere  *baseClause
	returning      []Column
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

func (c *insertClause) OnConflict(columns ...Column) *insertClause {
	c.onConflict = true
	c.conflictTarget = columns
	return c
}

func (c *insertClause) OnConstraint(constraint Node) *insertClause {
	c.onConstraint = constraint
	return c
}

func (c *insertClause) DoNothing() *insertClause {
	return c
}

func (c *insertClause) DoUpdateSet(setMap map[Column]interface{}) *insertClause {
	c.doUpdateSet = setMap
	return c
}

func (c *insertClause) Where(n ...Node) *insertClause {
	c.doUpdateWhere = newBaseClause("WHERE", n...).setSplit(" AND ")
	return c
}

func (c *insertClause) Returning(columns ...Column) *insertClause {
	c.returning = columns
	return c
}

func (c *insertClause) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString("INSERT INTO ")
	sb.WriteString(c.tableName)

	if len(c.columns) > 0 {
		sb.WriteString(" ")
		writeColumns(sb, c.columns, true)
	}
	if c.values != nil {
		sb.WriteString(" ")
		sb.WriteString(c.values.SqlString())
	}
	if c.onConflict {
		sb.WriteString(" ON CONFLICT")
		if len(c.conflictTarget) > 0 {
			sb.WriteRune(' ')
			writeColumns(sb, c.conflictTarget, true)
		}
		if c.onConstraint != nil {
			sb.WriteString(" ON CONSTRAINT ")
			sb.WriteString(c.onConstraint.SqlString())
		}
		if len(c.doUpdateSet) > 0 {
			sb.WriteString(" DO UPDATE SET ")
			sets := make([]string, 0, len(c.doUpdateSet))
			for k, v := range c.doUpdateSet {
				set := k.ColumnName()
				set += " = "
				switch vt := v.(type) {
				case Column:
					set += "EXCLUDED." + vt.ColumnName()
				case Node:
					set += vt.SqlString()
				default:
					set += "%s"
				}
				sets = append(sets, set)
			}
			sb.WriteString(strings.Join(sets, ", "))
			if c.doUpdateWhere != nil && len(c.doUpdateWhere.nodes) > 0 {
				sb.WriteRune(' ')
				sb.WriteString(c.doUpdateWhere.SqlString())
			}
		} else {
			sb.WriteString(" DO NOTHING")
		}
	}
	if len(c.returning) > 0 {
		sb.WriteString(" RETURNING ")
		writeColumns(sb, c.returning, false)
	}

	return sb.String()
}

func (c *insertClause) Values() []interface{} {
	var vs []interface{}
	if c.values != nil {
		vs = append(vs, c.values.Values()...)
	}
	if c.onConflict && len(c.doUpdateSet) > 0 {
		for _, v := range c.doUpdateSet {
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
	return vs
}

func (c *insertClause) Query() *sqlf.Query {
	return sqlf.Sprintf(c.SqlString(), c.Values()...)
}
