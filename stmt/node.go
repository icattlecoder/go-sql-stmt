package stmt

import (
	"fmt"
)

type node interface {
	SqlString() string
}

type valueNode interface {
	Values() []interface{}
}

type underlyingNode interface {
	underlyingNode() node
}

func notNilNodes(nodes []node) []node {
	ns := make([]node, 0, len(nodes))
	for _, n := range nodes {
		if u, ok := n.(underlyingNode); ok {
			n = u.underlyingNode()
		}
		if n == nil {
			continue
		}
		ns = append(ns, n)
	}
	return ns
}

func getValues(nodes ...node) *value {
	v := value{}
	for _, n := range notNilNodes(nodes) {
		if u, ok := n.(underlyingNode); ok {
			n = u.underlyingNode()
		}
		if n == nil {
			continue
		}
		if vn, ok := n.(valueNode); ok && vn != nil {
			v.values = append(v.values, vn.Values()...)
		}
	}
	return &v
}

type value struct {
	values []interface{}
}

func (v *value) Values() []interface{} {
	return v.values
}

type basicTypeValue struct {
	value interface{}
	hold  string
}

//TODO: refactor with go2 generic type is a good option
func newBasicTypeValue(v interface{}) *basicTypeValue {
	switch v.(type) {
	case bool, uint8, uint16, uint32, uint64, int8, int16, int32, int, int64, float32, float64, string:
		return &basicTypeValue{value: v}
	default:
		//not support, do nothing
	}
	return &basicTypeValue{}
}

func (b *basicTypeValue) SqlString() string {
	return "%s"
}

func (b *basicTypeValue) Values() []interface{} {
	return []interface{}{b.value}
}

type operator struct {
	node
}

func (o *operator) Values() []interface{} {
	if v, ok := o.node.(valueNode); ok {
		return v.Values()
	}
	return nil
}

type alias struct {
	operator
	bracket bool
	label   string
}

func (a *alias) SqlString() string {
	sb := get()
	defer put(sb)

	if a.bracket {
		sb.WriteString("(")
	}
	sb.WriteString(a.node.SqlString())
	if a.bracket {
		sb.WriteString(")")
	}
	sb.WriteString(" AS ")
	sb.WriteString(a.label)
	return sb.String()
}

func newAlias(n string, bracket bool, node node) *alias {
	a := &alias{
		label:   n,
		bracket: bracket,
	}
	a.node = node
	return a
}

type args interface {
	node
}

type Bool bool

const (
	True  = Bool(true)
	False = Bool(false)
)

func (s Bool) SqlString() string {
	if s {
		return "TRUE"
	}
	return "FALSE"
}

type Int int64

func (s Int) SqlString() string {
	return fmt.Sprintf("%d", s)
}

type String string

func (s String) SqlString() string {
	return "%s"
}

func (s String) Values() []interface{} {
	return []interface{}{string(s)}
}

type rawString string

const (
	Null = rawString("NULL")
)

func (s rawString) SqlString() string {
	return string(s)
}

type branch struct {
	b bool
	node
}

func (b *branch) underlyingNode() node {
	return b.node
}

func If(t bool, node node) *branch {
	if t {
		return &branch{t, node}
	}
	return &branch{}
}

func (b *branch) Else(node node) node {
	return node
}

func (b *branch) ElseIf(t bool, node node) *branch {
	if t {
		b.b = true
		b.node = node
	}
	return b
}

func Nodes(n ...node) []node {
	return n
}
