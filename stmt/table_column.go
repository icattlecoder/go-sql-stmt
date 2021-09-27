package stmt

import "strings"

type star string

const All = star("")

func (s star) SqlString() string {
	return "*"
}

type Column string

func (t Column) SqlString() string {
	return string(t)
}

func (t Column) As(s string) *alias {
	return newAlias(s, false, t)
}

func (t Column) Eq(r node) node {
	return Equals(t, r)
}

//JsonField get json array element as text,
//注意f可能会被注入，需要调用者保证安全性
func (t Column) JsonField(f string) node {
	return &jsonOperator{
		c:  t,
		op: "->>",
		r:  String("'" + f + "'"),
	}
}

func (t Column) JsonContainsAny(f ...string) node {
	return &jsonOperator{
		c:  t,
		op: " @> ",
		r:  String("{" + strings.Join(f, ",") + "}"),
	}
}

func (t Column) JsonContainsAll(f ...string) node {
	return &jsonOperator{
		c:  t,
		op: " && ",
		r:  String("{" + strings.Join(f, ",") + "}"),
	}
}

func (t Column) EqString(r string) node {
	return Equals(t, newBasicTypeValue(r))
}

func (t Column) EqInt(val int) node {
	return Equals(t, newBasicTypeValue(val))
}

func (t Column) NotNull() node {
	return &comparisonOperator{op: " NOT NULL", l: t, r: nil}
}

func (t Column) IsNull() node {
	return &comparisonOperator{op: " IS NULL", l: t, r: nil}
}

func (t Column) Lt(r node) node {
	return LessThan(t, r)
}

func (t Column) LtInt(val int) node {
	return LessThan(t, newBasicTypeValue(val))
}

func (t Column) Lte(r node) node {
	return LessEquals(t, r)
}

func (t Column) LteInt(val int) node {
	return LessEquals(t, newBasicTypeValue(val))
}

func (t Column) Gt(r node) node {
	return GreaterThan(t, r)
}

func (t Column) GtInt(val int) node {
	return GreaterThan(t, newBasicTypeValue(val))
}

func (t Column) Gte(r node) node {
	return GreaterEquals(t, r)
}

func (t Column) GteInt(val int) node {

	return GreaterEquals(t, newBasicTypeValue(val))
}

func (t Column) Like(val string) node {
	return &comparisonOperator{op: "LIKE", l: t, r: newBasicTypeValue("%" + val + "%")}
}

func (t Column) LikePrefix(val string) node {

	return &comparisonOperator{op: "LIKE", l: t, r: newBasicTypeValue("%" + val)}
}

func (t Column) LikeSuffix(val string) node {

	return &comparisonOperator{op: "LIKE", l: t, r: newBasicTypeValue(val + "%")}
}

func (t Column) ILike(val string) node {

	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue("%" + val + "%")}
}

func (t Column) ILikePrefix(val string) node {

	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue("%" + val)}
}

func (t Column) ILikeSuffix(val string) node {
	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue(val + "%")}
}

func (t Column) Desc() node {
	return Desc(t)
}

func (t Column) Plus(n node) *arithmeticOperator {
	return &arithmeticOperator{
		source: t,
		items: []*arithmeticOperator{{
			op:     "+",
			source: n,
		}},
	}
}

func (t Column) Multi(n node) *arithmeticOperator {
	return &arithmeticOperator{
		source: t,
		items: []*arithmeticOperator{{
			op:     "*",
			source: n,
		}},
	}
}

//func (t Column) In(value node) node{
//	return &columnOp{l: t, op: "IN", r: &stmt.bracket{valueStmt: value}}
//}
