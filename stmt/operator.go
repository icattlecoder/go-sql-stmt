package stmt

type jsonOperator struct {
	c  Column
	op string
	r  Node
}

func (j *jsonOperator) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString(j.c.SqlString())
	sb.WriteString(j.op)
	sb.WriteString(j.r.SqlString())
	return sb.String()
}

func (j *jsonOperator) Values() []interface{} {
	return getValues(j.r).Values()
}
