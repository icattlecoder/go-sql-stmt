package stmt

type joinPredicate struct {
	*baseClause

	typ    string
	source Node
	on     []Node
}

func (j *joinPredicate) Values() []interface{} {
	valNodes := []Node{j.source}
	valNodes = append(valNodes, j.on...)
	return getValues(valNodes...).Values()
}

func (j *joinPredicate) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString(j.typ)
	sb.WriteString(" ")
	sb.WriteString(j.source.SqlString())

	if len(j.on) > 0 {
		sb.WriteString(" ON (")
	}
	for i, a := range j.on {
		sb.WriteString(a.SqlString())
		if i < len(j.on)-1 {
			sb.WriteString(" AND ")
		}
	}
	if len(j.on) > 0 {
		sb.WriteString(")")
	}
	return sb.String()
}

func (j *joinPredicate) On(cond ...Node) *joinPredicate {
	j.on = cond
	return j
}

func LeftJoin(source Node) *joinPredicate {

	return &joinPredicate{
		typ:    "LEFT JOIN",
		source: source,
	}
}

func Join(source Node) *joinPredicate {

	return &joinPredicate{
		typ:    "JOIN",
		source: source,
	}
}
