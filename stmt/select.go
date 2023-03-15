package stmt

import (
	"fmt"
	"strings"

	"github.com/keegancsmith/sqlf"
)

type selectClause struct {
	explain *explain
	_select *baseClause
	from    *from
	where   *baseClause
	groupBy *baseClause
	having  *baseClause
	orderBy *baseClause
	limit   *limitOffset
	offset  *limitOffset
	sub     string
}

func Select(n ...Node) *selectClause {
	return &selectClause{
		_select: newBaseClause("SELECT", n...),
	}
}

// SubQuery 将当前查询作为子查询
func (c *selectClause) SubQuery(n string) *selectClause {
	c.sub = n
	return c
}

func (c *selectClause) From(sources ...Node) *selectClause {
	sources = notNilNodes(sources)
	if len(sources) == 0 {
		return c
	}
	var js []*joinPredicate
	if len(sources) > 1 {
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
	}
	f := from{
		sources: sources[0],
		joins:   js,
	}
	c.from = &f
	return c
}

func (c *selectClause) Where(n ...Node) *selectClause {
	c.where = newBaseClause("WHERE", n...).setSplit(" AND ")
	return c
}

func (c *selectClause) OrderBy(n ...Node) *selectClause {
	c.orderBy = newBaseClause("ORDER BY", n...)
	return c
}

func (c *selectClause) GroupBy(n ...Node) *selectClause {
	c.groupBy = newBaseClause("GROUP BY", n...)
	return c
}

func (c *selectClause) Having(n ...Node) *selectClause {
	c.having = newBaseClause("HAVING BY", n...)
	return c
}

func (c *selectClause) PartitionBy(n ...Node) *selectClause {
	c.having = newBaseClause("PARTITION BY", n...)
	return c
}

func (c *selectClause) Limit(size int) *selectClause {
	c.limit = &limitOffset{
		label: "LIMIT",
		val:   size,
	}
	return c
}

func (c *selectClause) ExplainJson() *selectClause {
	nc := *c
	nc.explain = &explain{Explain{Format: "JSON"}}
	return &nc
}

func (c *selectClause) Distinct() *selectClause {
	c._select.label = "SELECT DISTINCT"
	return c
}

func (c *selectClause) Offset(skip int) *selectClause {
	c.offset = &limitOffset{
		label: "OFFSET",
		val:   skip,
	}
	return c
}

func (c *selectClause) As(n string) *alias {
	return newAlias(n, true, c)
}

func (c *selectClause) SqlString() string {
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

	if c.sub != "" {
		return fmt.Sprintf("(%s) %s", sb.String(), c.sub)

	}
	return sb.String()
}

func (c *selectClause) Values() []interface{} {

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

func (c *selectClause) Union(c2 *selectClause) *setOperator {

	setOp := &setOperator{
		operator: "UNION",
		source:   &bracket{c},
	}
	setOp.Union(c2)
	setOp.selectClause = &selectClause{
		_select: newBaseClause("", setOp),
	}
	return setOp
}

func (c *selectClause) Query() *sqlf.Query {
	return sqlf.Sprintf(c.SqlString(), c.Values()...)
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
	case *selectClause:
		return &bracket{Node: n}
	default:
		return n
	}
}

func (b *bracket) SqlString() string {
	return "(" + b.Node.SqlString() + ")"
}
