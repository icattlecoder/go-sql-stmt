package stmt

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
