package stmt

import (
	"encoding/json"
	"reflect"
	"strings"
)

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

func (t Column) Eq(r Node) Node {
	return Equals(t, r)
}

//JsonField get json array element as text,
//注意f可能会被注入，需要调用者保证安全性
func (t Column) JsonField(f string) Node {
	return &jsonOperator{
		c:  t,
		op: "->>",
		r:  rawString("'" + f + "'"),
	}
}

//JsonContainsKevValue 用于过滤json数组里面的kv
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

func (t Column) JsonContainsAny(f ...string) Node {
	return &jsonOperator{
		c:  t,
		op: " @> ",
		r:  String("{" + strings.Join(f, ",") + "}"),
	}
}

func (t Column) JsonContainsAll(f ...string) Node {
	return &jsonOperator{
		c:  t,
		op: " && ",
		r:  String("{" + strings.Join(f, ",") + "}"),
	}
}

func (t Column) EqString(r string) Node {
	return Equals(t, newBasicTypeValue(r))
}

func (t Column) EqInt(val int) Node {
	return Equals(t, newBasicTypeValue(val))
}

func (t Column) NotNull() Node {
	return &comparisonOperator{op: " NOT NULL", l: t, r: nil}
}

func (t Column) IsNull() Node {
	return &comparisonOperator{op: " IS NULL", l: t, r: nil}
}

func (t Column) Lt(r Node) Node {
	return LessThan(t, r)
}

func (t Column) LtInt(val int) Node {
	return LessThan(t, newBasicTypeValue(val))
}

func (t Column) Lte(r Node) Node {
	return LessEquals(t, r)
}

func (t Column) LteInt(val int) Node {
	return LessEquals(t, newBasicTypeValue(val))
}

func (t Column) Gt(r Node) Node {
	return GreaterThan(t, r)
}

func (t Column) GtInt(val int) Node {
	return GreaterThan(t, newBasicTypeValue(val))
}

func (t Column) Gte(r Node) Node {
	return GreaterEquals(t, r)
}

func (t Column) GteInt(val int) Node {

	return GreaterEquals(t, newBasicTypeValue(val))
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

func (t Column) ILike(val string) Node {

	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue("%" + val + "%")}
}

func (t Column) ILikePrefix(val string) Node {

	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue("%" + val)}
}

func (t Column) ILikeSuffix(val string) Node {
	return &comparisonOperator{op: "ILIKE", l: t, r: newBasicTypeValue(val + "%")}
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

func (t Column) InString(val ...string) Node {
	return &comparisonOperator{op: "IN", l: t, r: expandRaw(val)}
}

func (t Column) In(n Node) Node {
	return &comparisonOperator{op: "IN", l: t, r: ibracket(n)}
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

//func (t Column) In(value Node) Node{
//	return &columnOp{l: t, op: "IN", r: &stmt.bracket{valueStmt: value}}
//}