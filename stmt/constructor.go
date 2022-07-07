package stmt

import (
	"strings"
)

type constructor struct {
	Name   string
	values []interface{}
}

type arrayConstructor struct {
	constructor
}

func (a *arrayConstructor) Values() []interface{} {
	var vs []interface{}
	for _, v := range a.values {
		vs = append(vs, v)
	}
	return vs
}

func (a *arrayConstructor) SqlString() string {
	sb := get()
	defer put(sb)
	sb.WriteString(a.Name)
	sb.WriteString("[")
	valuesSqlString(sb, a.values...)
	sb.WriteString("]")
	return sb.String()
}

func valuesSqlString(sb *strings.Builder, args ...interface{}) {
	for i := range args {
		sb.WriteString("%s")
		if i < len(args)-1 {
			sb.WriteString(", ")
		}
	}
}

func Array(values ...interface{}) *arrayConstructor {
	return &arrayConstructor{
		constructor{
			Name:   "ARRAY",
			values: values,
		},
	}
}
