package stmt

import (
	"fmt"
	"time"
)

type Node interface {
	SqlString() string
}

type valueNode interface {
	Values() []interface{}
}

type underlyingNode interface {
	underlyingNode() Node
}

func notNilNodes(nodes []Node) []Node {
	ns := make([]Node, 0, len(nodes))
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

func getValues(nodes ...Node) *value {
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
	case []int, []string:
		return &basicTypeValue{value: v}
	case time.Time:
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
	Node
}

func (o *operator) Values() []interface{} {
	if v, ok := o.Node.(valueNode); ok {
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
	sb.WriteString(a.Node.SqlString())
	if a.bracket {
		sb.WriteString(")")
	}
	sb.WriteString(" AS ")
	sb.WriteString(a.label)
	return sb.String()
}

func newAlias(n string, bracket bool, node Node) *alias {
	a := &alias{
		label:   n,
		bracket: bracket,
	}
	a.Node = node
	return a
}

type args interface {
	Node
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

type rawQuery struct {
	q string
	v []interface{}
}

func (r *rawQuery) SqlString() string {
	return r.q
}

func (r *rawQuery) Values() []interface{} {
	return r.v
}

type branch struct {
	b bool
	Node
}

func (b *branch) underlyingNode() Node {
	return b.Node
}

func If(t bool, node Node) *branch {
	if t {
		return &branch{t, node}
	}
	return &branch{}
}

func (b *branch) Else(node Node) Node {
	return node
}

func (b *branch) ElseIf(t bool, node Node) *branch {
	if t {
		b.b = true
		b.Node = node
	}
	return b
}

func Nodes(n ...Node) []Node {
	return n
}
