# go-sql-stmt

## Example

`stmt/node_test.go`

## 定义 Schema

可以根据json格式的schema定义来生成go struct

```javascript
//schema.json
[{
    "table": "channels",
    "columns": [
        "id",
        "uid",
        "name",
        "created_at"
    ]
}]
```

如果schema.json在当前路径，直接运行schema-generator

```go
//go:generate schema-generator
```

生成文件schema.go文件

```go
package db

import "github.com/fork-ai/go-sql-stmt/stmt"

var (
	// Channels is table channels
	Channels = newChannels("channels", "")
)

type channels struct {
	table, alias string
	Id           stmt.Column //Id is column of id
	Uid          stmt.Column //Uid is column of uid
	Name         stmt.Column //Name is column of name
	CreatedAt    stmt.Column //CreatedAt is column of created_at

}

func newChannels(table, alias string) channels {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return channels{
		table:     table,
		alias:     alias,
		Id:        stmt.Column(prefix + ".id"),
		Uid:       stmt.Column(prefix + ".uid"),
		Name:      stmt.Column(prefix + ".name"),
		CreatedAt: stmt.Column(prefix + ".created_at"),
	}
}
```

## 为什么使用Go构造sql？

1. 通过变量定义表和列的名称，只需要定义schema，避免输入错误的字符串引起调试困难。 例如有下面的sql:

```sql
SELECT download_counts.product_id,
       download_counts.year,
       SUM(download_counts.count) AS total_counts
FROM download_counts
WHERE download_counts.channel_id = %s
ORDER BY download_counts.product_id, download_counts.year DESC
    LIMIT 20
OFFSET 100
```

如果使用sql-stmt进行构造，代码如下:

```go
package main

import (
	. "github.com/fork-ai/go-sql-stmt/stmt"
)

func main() {
	Select(
		DownloadCounts.ProductId,
		DownloadCounts.Year,
		Sum(DownloadCounts.Count).As("total_counts"),
	).
		From(DownloadCounts).
		Where(DownloadCounts.ChannelId.EqInt(1)).
		GroupBy(DownloadCounts.ProductId, DownloadCounts.Year).
		OrderBy(DownloadCounts.Year.Desc()).
		Limit(20).
		Offset(100)
}
```

2. 明确的变量引用，不需要通过Printf来手动对齐占位符与变量。 例如上例中，ChannelId添加了 `channel_id=1`的限制。
3. 支持条件语句 有时我们会根据用户的输入值来添加Where子句中的限定，如果使用sql字符串拼接，则整个sql会被切割打散，这样的代码不易维护,例如

```go
package main

import (
	"github.com/keegancsmith/sqlf"
)

func main() {

	conds := []*sqlf.Query{
		sqlf.Sprintf("TRUE"),
	}
	switch opts.Type {
	case RankTypeFree:
		conds = append(conds, sqlf.Sprintf(`paid = FALSE`))
	case RankTypePaid:
		conds = append(conds, sqlf.Sprintf(`paid = TRUE`))
	}
	if opts.Category != "" {
		conds = append(conds, sqlf.Sprintf(`tags @> %s`, "{"+opts.Category+"}"))
	}
	if opts.ChannelUID == "apple" && opts.Device == AppleDeviceIPad {
		conds = append(conds, sqlf.Sprintf(`products_stats.extra->>'ipad' = 'true'`))
	}

	_ = sqlf.Sprintf(`
SELECT
	ROW_NUMBER() OVER(ORDER BY SUM(downloads) DESC) AS order,
	id, name, logo_url, uid, channel_uid, publisher_name, publisher_uid, SUM(downloads) AS download_count
FROM products_stats
WHERE
	channel_uid = %s AND country = %s
AND %s
GROUP BY id, name, logo_url, uid, channel_uid, publisher_name, publisher_uid
LIMIT %d OFFSET %d`)

}
```

下载使用sql-stmt来重构这段代码

```go
package main

import (
	. "github.com/fork-ai/go-sql-stmt/stmt"
)

type ProductsRankOptions struct {
	Type                 string
	ChannelUID           string
	Country              string
	Category             string
	Device               string // 仅针对 apple 有效，可以是 iphone 或 ipad，默认为 iphone
	IncludeDownloadCount bool
	Limit                int
	Offset               int
}

func main() {

	rankOpt := ProductsRankOptions{}
	q := Select(
		RowNumber().Over(OrderBy(Desc(Sum(ProductsStats.Downloads)))).As("order"),
		ProductsStats.Id,
		ProductsStats.Name,
		ProductsStats.LogoUrl,
		ProductsStats.Uid,
		ProductsStats.ChannelUid,
		ProductsStats.PublisherName,
		ProductsStats.PublisherUid,
		Sum(ProductsStats.Downloads).As("download_count"),
	).
		From(
			ProductsStats,
		).
		Where(
			If(rankOpt.ChannelUID != "", ProductsStats.ChannelUid.EqString(rankOpt.ChannelUID)),
			If(rankOpt.Type == "free",
				ProductsStats.Paid.Eq(False),
			).ElseIf(rankOpt.Type == "paid", ProductsStats.Paid.Eq(True)),
			If(rankOpt.ChannelUID == "apple" && rankOpt.Device == "ipad",
				Equals(ProductsStats.Extra.JsonField("ipad"), String("true")),
			),
		).
		GroupBy(
			ProductsStats.Id,
			ProductsStats.Name,
			ProductsStats.LogoUrl,
			ProductsStats.Uid,
			ProductsStats.ChannelUid,
			ProductsStats.PublisherName,
			ProductsStats.PublisherUid,
		).
		Limit(10)
	_ = q.Query()
}
```

`If`语句的一般使用方法是

```go
If(c1, node1).
    ElseIf(c2, node2).
    ElseIf(c3, node3).
    Else()
```

## 构造出来的sql是安全的吗？

字符串是引起sql注入的原因。因此在处理字符串时，我们会使用变量占位的方式。而其它类型如整数则在在构造sql时就直接生成了。理论上来讲，是安全的。