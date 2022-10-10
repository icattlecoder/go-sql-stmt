package stmt

import (
	"strings"
)

type arrayConstructor struct {
	values []interface{}
}

func Array(values ...interface{}) *arrayConstructor {
	return &arrayConstructor{
		values: values,
	}
}

func (c *arrayConstructor) SqlString() string {
	if len(c.values) == 0 {
		return ""
	}
	sb := get()
	defer put(sb)
	sb.WriteString("ARRAY")
	c.writeValues(sb, c.values)
	return sb.String()
}

func (c *arrayConstructor) Values() []interface{} {
	return c.getValues(c.values)
}

func (c *arrayConstructor) writeValues(sb *strings.Builder, values []interface{}) {
	sb.WriteString("[")
	for i, v := range values {
		switch vt := v.(type) {
		case *arrayConstructor:
			c.writeValues(sb, vt.values)
		default:
			sb.WriteString("%s")
		}
		if i < len(values)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
}

func (c *arrayConstructor) getValues(src []interface{}) []interface{} {
	vs := make([]interface{}, 0)
	for _, v := range src {
		switch vt := v.(type) {
		case *arrayConstructor:
			vs = append(vs, c.getValues(vt.values)...)
		default:
			vs = append(vs, v)
		}
	}
	return vs
}

type valuesConstructor struct {
	values []interface{}
}

func Values(values ...interface{}) *valuesConstructor {
	return &valuesConstructor{
		values: values,
	}
}

func (c *valuesConstructor) SqlString() string {
	if len(c.values) == 0 {
		return ""
	}
	sb := get()
	defer put(sb)
	sb.WriteString("VALUES ")

	vcs := c.nestValues()

	for i, v := range vcs {
		sb.WriteString("(")
		for vi, vv := range v.values {
			switch vt := vv.(type) {
			case Node:
				sb.WriteString(vt.SqlString())
			default:
				sb.WriteString("%s")
			}
			if vi < len(v.values)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
		if i < len(vcs)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (c *valuesConstructor) Values() []interface{} {
	vs := make([]interface{}, 0)

	vcs := c.nestValues()

	for _, v := range vcs {
		for _, vv := range v.values {
			switch vt := vv.(type) {
			case valueNode:
				vs = append(vs, vt.Values()...)
			default:
				vs = append(vs, vv)
			}
		}
	}

	return vs
}

func (c *valuesConstructor) As(n string) *alias {
	return newAlias(n, true, c)
}

func (c *valuesConstructor) nestValues() []*valuesConstructor {
	vcs := make([]*valuesConstructor, 0, len(c.values))
	first := c.values[0]
	if first != nil {
		if _, ok := first.(*valuesConstructor); ok {
			for _, v := range c.values {
				if vc, vok := v.(*valuesConstructor); vok {
					vcs = append(vcs, vc)
				}
			}
		} else {
			vcs = append(vcs, &valuesConstructor{
				values: c.values,
			})
		}
	}
	return vcs
}
