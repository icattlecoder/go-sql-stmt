package stmt

import (
	"fmt"
	"strings"

	"github.com/keegancsmith/sqlf"
)

type baseClause struct {
	label string
	nodes []Node
	split string
	valueNode
}

func (b *baseClause) Values() []interface{} {
	return getValues(b.nodes...).Values()
}

func newBaseClause(l string, nodes ...Node) *baseClause {
	nodes = notNilNodes(nodes)
	return &baseClause{
		label:     l,
		nodes:     nodes,
		split:     ", ",
		valueNode: getValues(nodes...),
	}
}

func (b *baseClause) setSplit(s string) *baseClause {
	b.split = s
	return b
}

func (b *baseClause) SqlString() string {
	sb := get()
	defer put(sb)

	if b.label != "" {
		sb.WriteString(b.label)
		sb.WriteString(" ")
	}
	for i, n := range b.nodes {
		if n == nil {
			continue
		}
		sb.WriteString(n.SqlString())
		if i < len(b.nodes)-1 {
			sb.WriteString(b.split)
		}
	}
	return sb.String()
}

type clause struct {
	explain *explain
	_select *baseClause
	from    *from
	where   *baseClause
	groupBy *baseClause
	having  *baseClause
	orderBy *baseClause
	limit   *limitOffset
	offset  *limitOffset
}

type from struct {
	joins   []*joinPredicate
	sources Node
}

func (f *from) Values() []interface{} {
	vals := getValues(f.sources).Values()
	for _, j := range f.joins {
		vals = append(vals, getValues(j).Values()...)
	}
	return vals
}

func (f *from) SqlString() string {
	sb := get()
	defer put(sb)
	sb.WriteString("FROM ")
	sb.WriteString(f.sources.SqlString())
	if len(f.joins) > 0 {
		sb.WriteString(" ")
		for i, j := range f.joins {
			sb.WriteString(j.SqlString())
			if i < len(f.joins)-1 {
				sb.WriteString(" ")
			}
		}
	}
	return sb.String()
}

type limitOffset struct {
	label string
	val   int
}

func (l *limitOffset) SqlString() string {
	return l.label + fmt.Sprintf(" %d", l.val)
}

func Select(n ...Node) *clause {
	return &clause{
		_select: newBaseClause("SELECT", n...),
	}
}

func (c *clause) From(sources ...Node) *clause {

	var js []*joinPredicate
	//FIXME panic
	for _, j := range sources[1:] {
		if p, ok := j.(*joinPredicate); ok {
			js = append(js, p)
			continue
		}
		js = append(js, &joinPredicate{
			typ:    "CROSS JOIN",
			source: j,
		})
	}

	f := from{
		sources: sources[0],
		joins:   js,
	}
	c.from = &f
	return c
}

func (c *clause) Where(n ...Node) *clause {
	c.where = newBaseClause("WHERE", n...).setSplit(" AND ")
	return c
}

func (c *clause) OrderBy(n ...Node) *clause {
	c.orderBy = newBaseClause("ORDER BY", n...)
	return c
}

func (c *clause) GroupBy(n ...Node) *clause {
	c.groupBy = newBaseClause("GROUP BY", n...)
	return c
}

func (c *clause) Having(n ...Node) *clause {
	c.having = newBaseClause("HAVING BY", n...)
	return c
}

func (c *clause) PartitionBy(n ...Node) *clause {
	c.having = newBaseClause("PARTITION BY", n...)
	return c
}

func (c *clause) Limit(size int) *clause {
	c.limit = &limitOffset{
		label: "LIMIT",
		val:   size,
	}
	return c
}

func (c *clause) ExplainJson() *clause {
	nc := *c
	nc.explain = &explain{Explain{Format: "JSON"}}
	return &nc
}

func (c *clause) Distinct() *clause {
	c._select.label = "SELECT DISTINCT"
	return c
}

func (c *clause) Offset(skip int) *clause {
	c.offset = &limitOffset{
		label: "OFFSET",
		val:   skip,
	}
	return c
}

func (c *clause) As(n string) *alias {
	return newAlias(n, true, c)
}

func (c *clause) SqlString() string {
	sb := get()
	defer put(sb)

	var bs []Node
	if c.explain != nil {
		bs = append(bs, c.explain)
	}

	if c._select != nil {
		bs = append(bs, c._select)
	}
	if c.from != nil {
		bs = append(bs, c.from)
	}
	if c.where != nil && len(c.where.nodes) > 0 {
		bs = append(bs, c.where)
	}
	if c.groupBy != nil && len(c.groupBy.nodes) > 0 {
		bs = append(bs, c.groupBy)
	}
	if c.having != nil && len(c.having.nodes) > 0 {
		bs = append(bs, c.having)
	}
	if c.orderBy != nil && len(c.orderBy.nodes) > 0 {
		bs = append(bs, c.orderBy)
	}
	if c.limit != nil {
		bs = append(bs, c.limit)
	}
	if c.offset != nil {
		bs = append(bs, c.offset)
	}

	for i, b := range bs {
		sb.WriteString(b.SqlString())
		if i < len(bs)-1 {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

func (c *clause) Values() []interface{} {

	var vs []interface{}
	if c._select != nil {
		vs = append(vs, c._select.Values()...)
	}
	if c.from != nil {
		vs = append(vs, c.from.Values()...)
	}
	if c.where != nil {
		vs = append(vs, c.where.Values()...)
	}
	if c.groupBy != nil {
		vs = append(vs, c.groupBy.Values()...)
	}
	if c.having != nil {
		vs = append(vs, c.having.Values()...)
	}
	if c.orderBy != nil {
		vs = append(vs, c.orderBy.Values()...)
	}
	return vs
}

func (c *clause) Union(c2 *clause) *setOperator {

	setOp := &setOperator{
		operator: "UNION",
		source:   &bracket{c},
	}
	setOp.Union(c2)
	setOp.clause = &clause{
		_select: newBaseClause("", setOp),
	}
	return setOp
}

func (c *clause) Query() *sqlf.Query {
	return sqlf.Sprintf(c.SqlString(), c.Values()...)
}

type Explain struct {
	//TEXT, XML, JSON, YAML
	Format string
}

var validFormats = map[string]struct{}{
	"XML":  {},
	"JSON": {},
	"YAML": {},
}

type explain struct {
	Explain
}

func (e *explain) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString("EXPLAIN")

	var opts []string
	if _, ok := validFormats[e.Format]; ok {
		opts = append(opts, fmt.Sprintf("FORMAT %s", e.Format))
	}

	if len(opts) > 0 {
		sb.WriteString(" (")
		sb.WriteString(strings.Join(opts, ","))
		sb.WriteString(")")
	}
	return sb.String()
}

type orderDirection struct {
	dir string
	Node
}

func (o *orderDirection) SqlString() string {
	return o.Node.SqlString() + " " + o.dir
}

func Asc(n Node) *orderDirection {
	return &orderDirection{
		dir:  "ASC",
		Node: n,
	}
}

func Desc(n Node) *orderDirection {
	return &orderDirection{
		dir:  "DESC",
		Node: n,
	}
}

func SortDir(desc bool, n Node) *orderDirection {
	if desc {
		return Desc(n)
	}
	return Asc(n)
}

type not struct {
	Node
}

func Not(n Node) *not {
	return &not{Node: n}
}

func (n *not) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString("NOT(")
	sb.WriteString(n.Node.SqlString())
	sb.WriteString(")")
	return sb.String()
}

func (n *not) Values() []interface{} {
	return getValues(n.Node).Values()
}

type bracket struct {
	Node
}

func (b *bracket) Values() []interface{} {
	return getValues(b.Node).Values()
}

func ibracket(n Node) Node {
	switch n.(type) {
	case *clause:
		return &bracket{Node: n}
	default:
		return n
	}
}

func (b *bracket) SqlString() string {
	return "(" + b.Node.SqlString() + ")"
}
