package stmt

type arithmeticOperator struct {
	op     string
	source node
	items  []*arithmeticOperator
}

func (a *arithmeticOperator) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString(a.source.SqlString())
	for _, i := range a.items {
		sb.WriteString(i.op)
		sb.WriteString("(")
		sb.WriteString(i.source.SqlString())
		sb.WriteString(")")
	}
	return sb.String()
}

func (a *arithmeticOperator) Plus(n node) *arithmeticOperator {
	a.items = append(a.items, &arithmeticOperator{
		op:     "+",
		source: n,
	})
	return a
}

func (a *arithmeticOperator) Minus(n node) *arithmeticOperator {
	a.items = append(a.items, &arithmeticOperator{
		op:     "-",
		source: n,
	})
	return a
}

func (a *arithmeticOperator) Mult(n node) *arithmeticOperator {
	a.items = append(a.items, &arithmeticOperator{
		op:     "*",
		source: n,
	})
	return a
}

func (a *arithmeticOperator) Div(n node) *arithmeticOperator {
	a.items = append(a.items, &arithmeticOperator{
		op:     "/",
		source: n,
	})
	return a
}
