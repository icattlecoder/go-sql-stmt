package stmt

type joinPredicate struct {
	*baseClause

	typ    string
	source node
	on     []node
}

func (j *joinPredicate) Values() []interface{} {
	return getValues(j.on...).Values()
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

func (j *joinPredicate) On(cond ...node) *joinPredicate {
	j.on = cond
	return j
}

func LeftJoin(source node) *joinPredicate {

	return &joinPredicate{
		typ:    "LEFT JOIN",
		source: source,
	}
}

func Join(source node) *joinPredicate {

	return &joinPredicate{
		typ:    "JOIN",
		source: source,
	}
}

