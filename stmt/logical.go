package stmt

type logicalOperator struct {
	op    string
	items []Node
}

func (l *logicalOperator) SqlString() string {
	sb := get()
	defer put(sb)
	sb.WriteString("(")

	for idx, i := range l.items {
		sb.WriteString(i.SqlString())
		if idx < len(l.items)-1 {
			sb.WriteString(l.op)
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func (l *logicalOperator) Values() []interface{} {
	return getValues(l.items...).Values()
}

func And(items ...Node) *logicalOperator {
	return &logicalOperator{
		op:    " AND ",
		items: notNilNodes(items),
	}
}

func Or(items ...Node) *logicalOperator {
	return &logicalOperator{
		op:    " OR ",
		items: notNilNodes(items),
	}
}
