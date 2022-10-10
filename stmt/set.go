package stmt

type setOperator struct {
	*selectClause
	operator string
	source   Node
	sets     []*setOperator
}

func (s *setOperator) SqlString() string {

	sb := get()
	defer put(sb)
	sb.WriteString(s.source.SqlString())
	for _, set := range s.sets {
		sb.WriteString(set.operator)
		sb.WriteString(set.source.SqlString())
	}
	return sb.String()
}

func (s *setOperator) Union(source Node) *setOperator {

	s.sets = append(s.sets, &setOperator{
		operator: "UNION",
		source:   &bracket{source},
	})
	return s
}

func (s *setOperator) Values() []interface{} {
	n := []Node{s.source}
	for _, su := range s.sets {
		n = append(n, su.source)
	}
	return getValues(n...).Values()
}
