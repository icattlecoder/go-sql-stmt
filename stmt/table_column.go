package stmt

import (
	"encoding/json"
	"reflect"
	"strings"
)

const All = Column("*")

type Column string

func (t Column) SqlString() string {
	return string(t)
}

func (t Column) ColumnName() string {
	sp := strings.Split(string(t), ".")
	return sp[len(sp)-1]
}

func (t Column) As(s string) *alias {
	return newAlias(s, false, t)
}

func (t Column) Eq(r Node) Node {
	return Equals(t, r)
}

// JsonField get json array element as text,
// 注意f可能会被注入，需要调用者保证安全性
func (t Column) JsonField(f string) Node {
	return &jsonOperator{
		c:  t,
		op: "->>",
		r:  rawString("'" + f + "'"),
	}
}

// JsonContainsKevValue 用于过滤json数组里面的kv
// 例如对jsonb字段 [{"round": "A", "name":"na"},{"round": "B", "name":"nb"}] 进行检索：
// JsonContainsKeyValue(map[string]interface{}{"round":"A"})
func (t Column) JsonContainsKevValue(kv map[string]interface{}) Node {

	bs, _ := json.Marshal(kv)
	r := &rawQuery{q: "%s::jsonb", v: []interface{}{"[" + string(bs) + "]"}}
	return &jsonOperator{
		c:  t,
		op: " @> ",
		r:  r,
	}
}

func (t Column) JsonContainsAll(f ...string) Node {
	return &jsonOperator{
		c:  t,
		op: " @> ",
		r:  String("{" + strings.Join(f, ",") + "}"),
	}
}

// ArrayContainsAll 包含 f 中所有元素
func (t Column) ArrayContainsAll(f ...string) Node {
	return &jsonOperator{
		c:  t,
		op: " @> ",
		r:  String("{" + strings.Join(f, ",") + "}"),
	}
}

// ArrayContainsAny 包含 f 中任意元素（t 和 f 中有至少一个元素相同）
func (t Column) ArrayContainsAny(f ...string) Node {
	return &jsonOperator{
		c:  t,
		op: " && ",
		r:  String("{" + strings.Join(f, ",") + "}"),
	}
}

func (t Column) EqString(r string) Node {
	return Equals(t, newBasicTypeValue(r))
}

func (t Column) EqBool(b bool) Node {
	return Equals(t, newBasicTypeValue(b))
}

func (t Column) EqInt(val int) Node {
	return Equals(t, newBasicTypeValue(val))
}

func (t Column) EqInt64(val int64) Node {
	return Equals(t, newBasicTypeValue(val))
}

func (t Column) NotNull() Node {
	return &comparisonOperator{op: "NOT NULL", l: t, r: nil}
}

func (t Column) IsNotNull() Node {
	return &comparisonOperator{op: "IS NOT NULL", l: t, r: nil}
}

func (t Column) IsNull() Node {
	return &comparisonOperator{op: "IS NULL", l: t, r: nil}
}

func (t Column) Lt(r Node) Node {
	return LessThan(t, r)
}

func (t Column) LtInt(val int) Node {
	return LessThan(t, newBasicTypeValue(val))
}

func (t Column) LtString(r string) Node {
	return LessThan(t, newBasicTypeValue(r))
}

func (t Column) Lte(r Node) Node {
	return LessEquals(t, r)
}

func (t Column) LteInt(val int) Node {
	return LessEquals(t, newBasicTypeValue(val))
}

func (t Column) LteString(r string) Node {
	return LessEquals(t, newBasicTypeValue(r))
}

func (t Column) Gt(r Node) Node {
	return GreaterThan(t, r)
}

func (t Column) GtInt(val int) Node {
	return GreaterThan(t, newBasicTypeValue(val))
}

func (t Column) GtString(r string) Node {
	return GreaterThan(t, newBasicTypeValue(r))
}

func (t Column) Gte(r Node) Node {
	return GreaterEquals(t, r)
}

func (t Column) GteInt(val int) Node {

	return GreaterEquals(t, newBasicTypeValue(val))
}

func (t Column) GteString(r string) Node {

	return GreaterEquals(t, newBasicTypeValue(r))
}

func (t Column) Like(val string) Node {
	return &comparisonOperator{op: "LIKE", l: t, r: newBasicTypeValue("%" + val + "%")}
}

func (t Column) LikePrefix(val string) Node {

	return &comparisonOperator{op: "LIKE", l: t, r: newBasicTypeValue("%" + val)}
}

func (t Column) LikeSuffix(val string) Node {

	return &comparisonOperator{op: "LIKE", l: t, r: newBasicTypeValue(val + "%")}
}

func (t Column) LikeRaw(val interface{}) Node {
	if n, ok := val.(Node); ok {
		return &comparisonOperator{op: "LIKE", l: t, r: n}
	}
	return &comparisonOperator{op: "LIKE", l: t, r: newBasicTypeValue(val)}
}

func (t Column) ILike(val string) Node {

	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue("%" + val + "%")}
}

func (t Column) ILikePrefix(val string) Node {

	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue("%" + val)}
}

func (t Column) ILikeSuffix(val string) Node {
	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue(val + "%")}
}

func (t Column) ILikeRaw(val interface{}) Node {
	if n, ok := val.(Node); ok {
		return &comparisonOperator{op: "ILIKE", l: t, r: n}
	}
	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue(val)}
}

// expand 会将指定数组对象平铺并重新拼接成用逗号分割的数据库查询参数。
func expandRaw(v interface{}) *rawQuery {
	kind := reflect.TypeOf(v).Kind()
	if kind != reflect.Slice {
		return nil
	}
	val := reflect.ValueOf(v)
	if val.Len() == 0 {
		return nil
	}
	q := rawQuery{}
	for i := 0; i < val.Len(); i++ {
		q.v = append(q.v, val.Index(i).Interface())
	}
	s := strings.Repeat("%s,", val.Len())
	q.q = "(" + s[:len(s)-1] + ")"
	return &q
}

func (t Column) InInt(val ...int) Node {
	return &comparisonOperator{op: "IN", l: t, r: expandRaw(val)}
}

func (t Column) InInt64(val ...int64) Node {
	return &comparisonOperator{op: "IN", l: t, r: expandRaw(val)}
}

func (t Column) InString(val ...string) Node {
	return &comparisonOperator{op: "IN", l: t, r: expandRaw(val)}
}

func (t Column) In(n Node) Node {
	return &comparisonOperator{op: "IN", l: t, r: ibracket(n)}
}

/*
Operator "@@" is Text Search Function indicates if tsvector matches tsquery
to_tsquery([ config regconfig , ] query text) : normalize words and convert to tsquery
*/
func (t Column) TsMatchText(regConfig, text string) Node {
	return &comparisonOperator{op: "@@ to_tsquery", l: t, r: expandRaw([]string{regConfig, text})}
}

/*
Operator "@@" is Text Search Function indicates if tsvector matches tsquery
plainto_tsquery([ config regconfig , ] query text) : tsquery produce tsquery ignoring punctuation
*/
func (t Column) TSMatchPlainText(regConfig, val string) Node {
	return &comparisonOperator{op: "@@ plainto_tsquery", l: t, r: expandRaw([]string{regConfig, val})}
}

func (t Column) Desc(isDesc ...bool) Node {
	if len(isDesc) == 0 {
		return Asc(t)
	}
	return SortDir(isDesc[0], t)
}

func (t Column) Plus(n Node) *arithmeticOperator {
	return &arithmeticOperator{
		source: t,
		items: []*arithmeticOperator{{
			op:     "+",
			source: n,
		}},
	}
}

func (t Column) Multi(n Node) *arithmeticOperator {
	return &arithmeticOperator{
		source: t,
		items: []*arithmeticOperator{{
			op:     "*",
			source: n,
		}},
	}
}

func (t Column) SimilarTo(val string) Node {
	return &comparisonOperator{op: "SIMILAR TO", l: t, r: newBasicTypeValue(val)}
}

func (t Column) Concatenates(v interface{}) Node {
	if n, ok := v.(Node); ok {
		return Concatenates(t, n)
	}
	return Concatenates(t, newBasicTypeValue(v))
}

func writeColumns(sb *strings.Builder, columns []Column, hasBracket bool) {
	if hasBracket {
		sb.WriteString("(")
	}
	for i, col := range columns {
		sb.WriteString(col.ColumnName())
		if i < len(columns)-1 {
			sb.WriteString(", ")
		}
	}
	if hasBracket {
		sb.WriteString(")")
	}
}
