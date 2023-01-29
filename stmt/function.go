package stmt

import "strings"

type function struct {
	Name string
	args []Node
}

func (f *function) Values() []interface{} {
	var vs []interface{}
	for _, a := range f.args {
		if v, ok := a.(valueNode); ok {
			vs = append(vs, v.Values()...)
		}
	}
	return vs
}

func (f *function) SqlString() string {

	sb := get()
	defer put(sb)
	sb.WriteString(f.Name)
	sb.WriteString("(")
	argsSqlString(sb, f.args...)
	sb.WriteString(")")
	return sb.String()
}

type scalarFunction struct {
	function
}

func Coalesce(x, y Node, others ...Node) *scalarFunction {

	return &scalarFunction{function: function{
		Name: "COALESCE",
		args: append([]Node{x, y}, others...),
	}}

}

type over struct {
	partition Node
	order     *baseClause
}

func (o *over) sqlString(sb *strings.Builder) string {

	sb.WriteString(" OVER (")
	if o.partition != nil {
		sb.WriteString("PARTITION BY " + o.partition.SqlString())
	}
	if o.order != nil {
		if o.partition != nil {
			sb.WriteString(" ")
		}
		sb.WriteString("ORDER BY " + o.order.SqlString())
	}
	sb.WriteString(")")
	return sb.String()
}

type windowFunction struct {
	function
	over *over
}

func Rank() *windowFunction {
	return &windowFunction{
		function: function{Name: "Rank"},
	}
}

func RowNumber() *windowFunction {
	return &windowFunction{
		function: function{Name: "ROW_NUMBER"},
	}
}

func DenseRank() *windowFunction {
	return &windowFunction{
		function: function{Name: "DENSE_RANK"},
	}
}

func FirstValue(n Node) *windowFunction {
	return &windowFunction{
		function: function{Name: "FIRST_VALUE", args: []Node{n}},
	}
}

func LastValue(n Node) *windowFunction {
	return &windowFunction{
		function: function{Name: "LAST_VALUE", args: []Node{n}},
	}
}

func NthValue(n Node, nth int) *windowFunction {
	return &windowFunction{
		function: function{Name: "NTH_VALUE", args: []Node{n, Int(nth)}},
	}
}

func (w *windowFunction) Over(n *over) *windowFunction {
	w.over = n
	return w
}

func (w *windowFunction) As(n string) *alias {
	return newAlias(n, false, w)
}

func PartitionBy(n Node) *over {

	return &over{
		partition: n,
	}
}

func OrderBy(n ...Node) *over {

	return &over{
		order: newBaseClause("", n...).setSplit(", "),
	}
}

func (o *over) OrderBy(n ...Node) *over {
	o.order = newBaseClause("", n...).setSplit(", ")
	return o
}

func (w *windowFunction) SqlString() string {
	sb := get()
	defer put(sb)
	sb.WriteString(w.function.SqlString())
	if w.over != nil {
		w.over.sqlString(sb)
	}
	return sb.String()
}

type aggregateFunction struct {
	function
}

func NewAggregateFunction(n string) func(args ...Node) *aggregateFunction {
	return func(args ...Node) *aggregateFunction {
		return &aggregateFunction{
			function{
				Name: n,
				args: notNilNodes(args),
			},
		}
	}
}

var (
	Sum         = NewAggregateFunction("SUM")
	Count       = NewAggregateFunction("COUNT")
	Avg         = NewAggregateFunction("AVG")
	Max         = NewAggregateFunction("MAX")
	Min         = NewAggregateFunction("MIN")
	ToJsonb     = NewAggregateFunction("to_jsonb")
	ArrayRemove = NewAggregateFunction("ARRAY_REMOVE")
	ArrayAgg    = NewAggregateFunction("ARRAY_AGG")
	Unnest      = NewAggregateFunction("UNNEST")
	Any         = NewAggregateFunction("ANY")
	Now         = NewAggregateFunction("NOW")
)

func (a *aggregateFunction) As(n string) *alias {
	return newAlias(n, false, a)
}

func (f *function) As(n string) *alias {
	return newAlias(n, false, f)
}

type distinct struct {
	args
}

func (d *distinct) SqlString() string {
	sb := get()
	defer put(sb)
	sb.WriteString("DISTINCT ")
	sb.WriteString(d.args.SqlString())
	return sb.String()
}

func Distinct(args args) *distinct {
	return &distinct{args: args}
}
