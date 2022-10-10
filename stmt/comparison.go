package stmt

type comparisonOperator struct {
	op   string
	l, r Node
}

func NewComparisonOperator(op string) func(l, r Node) *comparisonOperator {
	return func(l, r Node) *comparisonOperator {
		return &comparisonOperator{op: op, l: ibracket(l), r: ibracket(r)}
	}
}

func (c *comparisonOperator) SqlString() string {
	sb := get()
	defer put(sb)

	sb.WriteString(c.l.SqlString())
	sb.WriteRune(' ')
	sb.WriteString(c.op)
	if c.r != nil {
		sb.WriteRune(' ')
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
