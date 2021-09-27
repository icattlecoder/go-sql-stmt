package stmt

type comparisonOperator struct {
	op   string
	l, r node
}

func NewComparisonOperator(op string) func(l, r node) *comparisonOperator {
	return func(l, r node) *comparisonOperator {
		return &comparisonOperator{op: op, l: l, r: r}
	}
}

func (c *comparisonOperator) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString(c.l.SqlString())
	sb.WriteString(" " + c.op + " ")
	if c.r != nil {
		sb.WriteString(c.r.SqlString())
	}
	return sb.String()
}

func (c *comparisonOperator) Values() []interface{} {
	return getValues(c.l, c.r).Values()
}

var (
	GreaterThan   = NewComparisonOperator(">")
	LessThan      = NewComparisonOperator("<")
	Equals        = NewComparisonOperator("=")
	NotEquals     = NewComparisonOperator("!=")
	GreaterEquals = NewComparisonOperator(">=")
	LessEquals    = NewComparisonOperator("<=")
)
