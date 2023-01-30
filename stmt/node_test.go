package stmt

import (
	"reflect"
	"strings"
	"testing"
)

//go:generate schema-generator

func removeNewLine(s string) string {
	return strings.ReplaceAll(s, "\n", "")
}

func isStrEq(a, b string) bool {
	return removeNewLine(a) == removeNewLine(b)
}

func TestSelectClauses(t *testing.T) {

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

	rankOpt := ProductsRankOptions{
		Type:                 "free",
		ChannelUID:           "",
		Country:              "",
		Category:             "",
		Device:               "",
		IncludeDownloadCount: false,
		Limit:                0,
		Offset:               0,
	}

	var (
		keyword = "wechat"
	)
	_ = keyword

	tests := []struct {
		name       string
		clause     *selectClause
		wantQuery  string
		wantValues []interface{}
	}{
		{
			name: "simple",
			clause: Select(
				DownloadCounts.ProductId,
				DownloadCounts.Year,
				Sum(DownloadCounts.Count).As("total_counts"),
			).
				From(DownloadCounts).
				Where(DownloadCounts.ChannelId.EqString("123"), DownloadCounts.ProductId.EqInt64(9223372036854775807)).
				GroupBy(DownloadCounts.ProductId, DownloadCounts.Year).
				OrderBy(DownloadCounts.ProductId, Desc(DownloadCounts.Year)).
				Limit(20).
				Offset(100),
			wantQuery: `SELECT 
download_counts.product_id, 
download_counts.year, 
SUM(download_counts.count) AS total_counts 
FROM download_counts 
WHERE download_counts.channel_id = %s AND download_counts.product_id = %s 
GROUP BY download_counts.product_id, download_counts.year 
ORDER BY download_counts.product_id, download_counts.year DESC 
LIMIT 20 OFFSET 100`,
			wantValues: []interface{}{"123", int64(9223372036854775807)},
		},
		{
			name: "union set operator",
			clause: Select(Products.Name).From(Products).Where(Or(Products.Name.ILike(keyword+"0"), Products.Description.ILike(keyword+"1"))).Limit(10).
				Union(
					Select(Services.Name).From(Services).Where(Or(Services.Name.ILike(keyword+"2"), Services.Description.ILike(keyword+"3"))).Limit(10),
				).
				Union(
					Select(Publishers.Name).From(Publishers).Where(Or(Publishers.Name.ILike(keyword+"4"), Publishers.Description.ILike(keyword+"5"))).Limit(10),
				).
				Limit(10),
			wantQuery: `(SELECT products.name FROM products WHERE (products.name ILIKE %s OR products.description ILIKE %s) LIMIT 10)
UNION
(SELECT services.name FROM services WHERE (services.name ILIKE %s OR services.description ILIKE %s) LIMIT 10)
UNION
(SELECT publishers.name FROM publishers WHERE (publishers.name ILIKE %s OR publishers.description ILIKE %s) LIMIT 10) 
LIMIT 10`,
			wantValues: []interface{}{"%wechat0%", "%wechat1%", "%wechat2%", "%wechat3%", "%wechat4%", "%wechat5%"},
		},
		{
			name: "window function",
			clause: Select(
				DownloadCounts.CountryCode,
				Rank().Over(PartitionBy(DownloadCounts.CountryCode).OrderBy(DownloadCounts.Year)).As("year_rank"),
				DenseRank().Over(PartitionBy(DownloadCounts.CountryCode).OrderBy(DownloadCounts.Year)).As("year_dense_rank"),
				PercentileDisc(0.5).WithinGroup(DownloadCounts.Year.Desc(true)).As("year_fifty_percent"),
				FirstValue(DownloadCounts.Count).Over(PartitionBy(DownloadCounts.CountryCode).OrderBy(DownloadCounts.Year)),
				LastValue(DownloadCounts.Count).Over(PartitionBy(DownloadCounts.CountryCode).OrderBy(DownloadCounts.Year)),
				NthValue(DownloadCounts.Count, 2).Over(PartitionBy(DownloadCounts.CountryCode).OrderBy(DownloadCounts.Year)),
			).
				From(DownloadCounts).
				Where(DownloadCounts.ProductId.EqInt(123)),
			wantQuery: `SELECT 
download_counts.country_code, 
Rank() OVER (PARTITION BY download_counts.country_code ORDER BY download_counts.year) AS year_rank, 
DENSE_RANK() OVER (PARTITION BY download_counts.country_code ORDER BY download_counts.year) AS year_dense_rank, 
PERCENTILE_DISC(%s) WITHIN GROUP (ORDER BY download_counts.year DESC) AS year_fifty_percent, 
FIRST_VALUE(download_counts.count) OVER (PARTITION BY download_counts.country_code ORDER BY download_counts.year), 
LAST_VALUE(download_counts.count) OVER (PARTITION BY download_counts.country_code ORDER BY download_counts.year), 
NTH_VALUE(download_counts.count, 2) OVER (PARTITION BY download_counts.country_code ORDER BY download_counts.year) 
FROM download_counts 
WHERE download_counts.product_id = %s`,
			wantValues: []interface{}{0.5, 123},
		},
		{
			name: "scalar function",
			clause: Select(
				Coalesce(Max(DownloadCounts.Count), Int(0)).As("max_downloads"),
			).
				From(DownloadCounts).
				Where(
					DownloadCounts.ChannelId.EqString("1"),
					DownloadCounts.ProductId.EqInt(123),
					DownloadCounts.CountryCode.EqString("CN"),
				),
			wantQuery: `SELECT COALESCE(MAX(download_counts.count), 0) AS max_downloads 
FROM download_counts WHERE download_counts.channel_id = %s 
AND download_counts.product_id = %s 
AND download_counts.country_code = %s`,
			wantValues: []interface{}{"1", 123, "CN"},
		},
		{
			name: "explain",
			clause: Select(
				ProductsQueryStats.Name,
			).
				From(ProductsQueryStats).
				Where(
					ProductsQueryStats.Name.EqString("ANDROID"),
				).Distinct().ExplainJson(),
			wantQuery: `EXPLAIN (FORMAT JSON) SELECT DISTINCT products_query_stats.name 
FROM products_query_stats WHERE products_query_stats.name = %s`,
			wantValues: []interface{}{"ANDROID"},
		},
		{
			name: "join",
			clause: Select(
				Products.Id,
				ArrayRemove(ArrayAgg(Distinct(DownloadCounts.CountryCode)), Null).As("countries"),
			).
				From(
					Products,
					Join(Publishers).On(Publishers.Id.Eq(Products.PublisherId)),
					Join(Channels).On(Channels.Id.Eq(Products.ChannelId)),
					LeftJoin(DownloadCounts).On(Products.Id.Eq(DownloadCounts.ProductId)),
					LeftJoin(Services).On(Services.Homepage.Eq(Any(Products.ScreenshotUrls))),
				).GroupBy(Products.Id, Publishers.Id, Channels.Id),
			wantQuery: `SELECT 
products.id, 
ARRAY_REMOVE(ARRAY_AGG(DISTINCT download_counts.country_code), NULL) AS countries 
FROM products 
JOIN publishers ON (publishers.id = products.publisher_id) 
JOIN channels ON (channels.id = products.channel_id) 
LEFT JOIN download_counts ON (products.id = download_counts.product_id) 
LEFT JOIN services ON (services.homepage = ANY(products.screenshot_urls)) 
GROUP BY products.id, publishers.id, channels.id`,
		},
		{
			name: "json function",
			clause: Select(
				All,
			).
				From(
					Products,
				).Where(
				Products.ChannelId.EqString("mi"),
				Equals(Products.Extra.JsonField("track_id"), String("abc")),
			).Limit(100),
			wantQuery: `SELECT * 
FROM products 
WHERE products.channel_id = %s 
AND products.extra->>'track_id' = %s 
LIMIT 100`,
			wantValues: []interface{}{"mi", "abc"},
		},
		{
			name: "if statement",
			clause: Select(
				RowNumber().Over(OrderBy(Desc(Sum(ProductsStats.Downloads)), ProductsStats.Paid)).As("order"),
				ProductsStats.Id,
				ProductsStats.Name,
				ProductsStats.LogoUrl,
				ProductsStats.Uid,
				ProductsStats.ChannelUid,
				ProductsStats.PublisherName,
				ProductsStats.PublisherUid,
				Sum(ProductsStats.Downloads).As("downloads_count"),
			).
				From(
					ProductsStats,
				).
				Where(
					ProductsStats.ChannelUid.EqString("mi"),
					ProductsStats.Country.EqString("CN"),
					If(rankOpt.Type == "free", ProductsStats.Paid.Eq(False)),
					If(rankOpt.Type == "paid", ProductsStats.Paid.Eq(True)),
					If(rankOpt.ChannelUID == "apple" && rankOpt.Device == "ipad",
						Equals(ProductsStats.Extra.JsonField("ipad"), String("true")),
					),
					Equals(ProductsStats.Extra.JsonField("track_id"), String("abc")),
				).Limit(100),
			wantQuery: `SELECT 
ROW_NUMBER() OVER (ORDER BY SUM(products_stats.downloads) DESC, products_stats.paid) AS order, 
products_stats.id, 
products_stats.name, 
products_stats.logo_url, 
products_stats.uid, 
products_stats.channel_uid, 
products_stats.publisher_name, 
products_stats.publisher_uid, 
SUM(products_stats.downloads) AS downloads_count 
FROM products_stats 
WHERE 
products_stats.channel_uid = %s 
AND products_stats.country = %s 
AND products_stats.paid = FALSE 
AND products_stats.extra->>'track_id' = %s LIMIT 100`,
			wantValues: []interface{}{"mi", "CN", "abc"},
		},
		{
			name: "arithmetic function",
			clause: Select(
				DownloadCounts.Year,
				DownloadCounts.Month,
				DownloadCounts.Count,
			).
				From(DownloadCounts).
				Where(
					DownloadCounts.ChannelId.EqInt(2),
					DownloadCounts.ProductId.EqInt(74831),
					Equals(
						DownloadCounts.Year.Multi(Int(100)).Plus(DownloadCounts.Month),
						Select(Max(DownloadCounts.Year.Multi(Int(100)).Plus(DownloadCounts.Month))).
							From(DownloadCounts).
							Where(
								DownloadCounts.ChannelId.EqInt(2),
								DownloadCounts.ProductId.EqInt(74831),
							),
					),
				),
			wantQuery: `SELECT 
download_counts.year, 
download_counts.month, 
download_counts.count FROM download_counts 
WHERE download_counts.channel_id = %s 
AND download_counts.product_id = %s 
AND download_counts.year*(100)+(download_counts.month) = (SELECT MAX(download_counts.year*(100)+(download_counts.month)) FROM download_counts WHERE download_counts.channel_id = %s AND download_counts.product_id = %s)`,
			wantValues: []interface{}{2, 74831, 2, 74831},
		},
		{
			name: "in predicate",
			clause: Select(
				DownloadCounts.Year,
			).
				From(DownloadCounts).
				Where(
					DownloadCounts.ChannelId.InInt(1, 2, 3),
					DownloadCounts.Year.InString("2010", "2011"),
				),
			wantQuery: `SELECT 
download_counts.year 
FROM download_counts 
WHERE download_counts.channel_id IN (%s,%s,%s) AND download_counts.year IN (%s,%s)`,
			wantValues: []interface{}{1, 2, 3, "2010", "2011"},
		},
		{
			name: "not",
			clause: Select(
				DownloadCounts.Year,
			).
				From(DownloadCounts).
				Where(
					Not(DownloadCounts.Year.InString("2010", "2011")),
				),
			wantQuery: `SELECT 
download_counts.year 
FROM download_counts 
WHERE NOT(download_counts.year IN (%s,%s))`,
			wantValues: []interface{}{"2010", "2011"},
		},
		{
			name: "array operator",
			clause: Select(
				All,
			).
				From(
					Products,
				).Where(
				Products.ScreenshotUrls.ArrayContainsAny("https://qq.com/image/100"),
				Products.Countries.ArrayContainsAll("CN", "US"),
			).Limit(100),
			wantQuery: `SELECT * 
FROM products 
WHERE products.screenshot_urls && %s 
AND products.countries @> %s 
LIMIT 100`,
			wantValues: []interface{}{"{https://qq.com/image/100}", "{CN,US}"},
		},
		{
			name: "subqueries",
			clause: Select(
				All,
			).
				From(
					Services,
					Join(
						Select(
							ProductsQueryStats.Id,
							ProductsQueryStats.Name,
							ProductsQueryStats.ServicesUids,
						).
							From(ProductsQueryStats).
							Where(ProductsQueryStats.BundleId.EqString("com.tencent.xin")).
							As("t1"),
					).On(Services.Uid.Eq(Any(Column("t1.services_uids")))),
				).
				Where(Services.PublisherId.GtInt(1)).
				OrderBy(Services.UpdatedAt),
			wantQuery: `SELECT * 
FROM services 
JOIN (SELECT products_query_stats.id, products_query_stats.name, products_query_stats.services_uids 
FROM products_query_stats 
WHERE products_query_stats.bundle_id = %s) AS t1 
ON (services.uid = ANY(t1.services_uids)) 
WHERE services.publisher_id > %s 
ORDER BY services.updated_at`,
			wantValues: []interface{}{"com.tencent.xin", 1},
		},
		{
			name: "similar to",
			clause: Select(
				All,
			).
				From(
					Services,
				).
				Where(
					Services.Name.SimilarTo("%(腾讯|阿里)%"),
				),
			wantQuery: `SELECT * 
FROM services 
WHERE services.name SIMILAR TO %s`,
			wantValues: []interface{}{"%(腾讯|阿里)%"},
		},
		{
			name: "array constructor",
			clause: Select(
				Array(1, 2, "3", 4.5),
			),
			wantQuery:  `SELECT ARRAY[%s, %s, %s, %s]`,
			wantValues: []interface{}{1, 2, "3", 4.5},
		},
		{
			name: "array constructor nest",
			clause: Select(
				Array(Array(1, 2), Array(3, 4)),
			),
			wantQuery:  `SELECT ARRAY[[%s, %s], [%s, %s]]`,
			wantValues: []interface{}{1, 2, 3, 4},
		},
		{
			name: "like array",
			clause: Select(
				All,
			).
				From(
					Services,
				).
				Where(
					Services.Name.LikeRaw(Any(Array("%对象%", "%存储%"))),
					Services.Name.ILikeRaw(Any(Array("uClouD%", "linkV%"))),
				),
			wantQuery: `SELECT * 
FROM services 
WHERE services.name LIKE ANY(ARRAY[%s, %s]) 
AND services.name ILIKE ANY(ARRAY[%s, %s])`,
			wantValues: []interface{}{"%对象%", "%存储%", "uClouD%", "linkV%"},
		},
		{
			name:       "values constructor",
			clause:     Select(All).From(Values(1, 2, "3", 4.5).As("t1")),
			wantQuery:  `SELECT * FROM (VALUES (%s, %s, %s, %s)) AS t1`,
			wantValues: []interface{}{1, 2, "3", 4.5},
		},
		{
			name:       "values constructor nest",
			clause:     Select(All).From(Values(Values(1, 2), Values(3, 4)).As("t1")),
			wantQuery:  `SELECT * FROM (VALUES (%s, %s), (%s, %s)) AS t1`,
			wantValues: []interface{}{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery := tt.clause.SqlString()
			gotValues := tt.clause.Values()
			if !isStrEq(gotQuery, tt.wantQuery) {
				t.Errorf("SqlString() = %v, want \n%v", gotQuery, removeNewLine(tt.wantQuery))
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestInsertClauses(t *testing.T) {
	tests := []struct {
		name       string
		clause     *insertClause
		wantQuery  string
		wantValues []interface{}
	}{
		{
			name: "simple",
			clause: InsertInto(
				Channels,
				[]Column{
					Channels.Id,
					Channels.Uid,
					Channels.Name,
					Channels.CreatedAt,
				},
				Values(
					1,
					"uid1",
					"channel1",
					Now(),
				),
			),
			wantQuery:  `INSERT INTO channels (id, uid, name, created_at) VALUES (%s, %s, %s, NOW())`,
			wantValues: []interface{}{1, "uid1", "channel1"},
		},
		{
			name: "multi values",
			clause: InsertInto(
				Channels,
				[]Column{
					Channels.Id,
					Channels.Uid,
					Channels.Name,
					Channels.CreatedAt,
				},
				Values(
					Values(
						1,
						"uid1",
						"channel1",
						Now(),
					),
					Values(
						2,
						"uid2",
						"channel2",
						Now(),
					),
				),
			),
			wantQuery:  `INSERT INTO channels (id, uid, name, created_at) VALUES (%s, %s, %s, NOW()), (%s, %s, %s, NOW())`,
			wantValues: []interface{}{1, "uid1", "channel1", 2, "uid2", "channel2"},
		},
		{
			name: "select values",
			clause: InsertInto(
				Channels,
				[]Column{
					Channels.Uid,
					Channels.Name,
					Channels.CreatedAt,
				},
				Select(
					Products.ChannelId,
					Products.ChannelId,
					Now(),
				).
					From(Products).
					Where(Products.Uid.EqString("product_1")),
			),
			wantQuery: `INSERT INTO channels (uid, name, created_at) 
SELECT products.channel_id, products.channel_id, NOW() 
FROM products WHERE products.uid = %s`,
			wantValues: []interface{}{"product_1"},
		},
		{
			name: "on conflict do nothing",
			clause: InsertInto(
				Channels,
				[]Column{
					Channels.Id,
					Channels.Uid,
					Channels.Name,
					Channels.CreatedAt,
				},
				Values(
					1,
					"uid1",
					"channel1",
					Now(),
				),
			).
				OnConflict(Channels.Uid).
				DoNothing(),
			wantQuery:  `INSERT INTO channels (id, uid, name, created_at) VALUES (%s, %s, %s, NOW()) ON CONFLICT (uid) DO NOTHING`,
			wantValues: []interface{}{1, "uid1", "channel1"},
		},
		{
			name: "on conflict do update set",
			clause: InsertInto(
				Channels,
				[]Column{
					Channels.Id,
					Channels.Uid,
					Channels.Name,
					Channels.CreatedAt,
				},
				Values(
					1,
					"uid1",
					"channel_1",
					Now(),
				),
			).
				OnConflict(Channels.Uid).
				DoUpdateSet(map[Column]interface{}{
					Channels.Uid:       "uid2",
					Channels.Name:      Channels.Name,
					Channels.CreatedAt: Now(),
				}).
				Where(Channels.Name.IsNotNull()),
			wantQuery: `INSERT INTO channels (id, uid, name, created_at) 
VALUES (%s, %s, %s, NOW()) 
ON CONFLICT (uid) 
DO UPDATE SET uid = %s, name = EXCLUDED.name, created_at = NOW() 
WHERE channels.name IS NOT NULL`,
			wantValues: []interface{}{1, "uid1", "channel_1", "uid2"},
		},
		{
			name: "on constraint",
			clause: InsertInto(
				Channels,
				[]Column{
					Channels.Id,
					Channels.Uid,
					Channels.Name,
					Channels.CreatedAt,
				},
				Values(
					1,
					"uid1",
					"channel1",
					Now(),
				),
			).
				OnConflict().
				OnConstraint(Column("channels_uid_key")),
			wantQuery: `INSERT INTO channels (id, uid, name, created_at) 
VALUES (%s, %s, %s, NOW()) 
ON CONFLICT ON CONSTRAINT channels_uid_key DO NOTHING`,
			wantValues: []interface{}{1, "uid1", "channel1"},
		},
		{
			name: "returning",
			clause: InsertInto(
				Channels,
				[]Column{
					Channels.Id,
					Channels.Uid,
					Channels.Name,
					Channels.CreatedAt,
				},
				Values(
					1,
					"uid1",
					"channel1",
					Now(),
				),
			).
				Returning(All),
			wantQuery:  `INSERT INTO channels (id, uid, name, created_at) VALUES (%s, %s, %s, NOW()) RETURNING *`,
			wantValues: []interface{}{1, "uid1", "channel1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery := tt.clause.SqlString()
			gotValues := tt.clause.Values()
			if !isStrEq(gotQuery, tt.wantQuery) {
				t.Errorf("SqlString() = %v, want \n%v", gotQuery, removeNewLine(tt.wantQuery))
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestUpdateClauses(t *testing.T) {
	tests := []struct {
		name       string
		clause     *updateClause
		wantQuery  string
		wantValues []interface{}
	}{
		{
			name: "simple",
			clause: Update(
				Channels,
			).
				Set(map[Column]interface{}{
					Channels.Name: "channel2",
				}).
				Where(Channels.Id.EqInt(1)),
			wantQuery:  `UPDATE channels SET name = %s WHERE channels.id = %s`,
			wantValues: []interface{}{"channel2", 1},
		},
		{
			name: "returning",
			clause: Update(
				Channels,
			).
				Set(map[Column]interface{}{
					Channels.Name: Channels.Uid,
				}).
				Where(Channels.Id.EqInt(1)).
				Returning(Channels.Uid),
			wantQuery:  `UPDATE channels SET name = uid WHERE channels.id = %s RETURNING uid`,
			wantValues: []interface{}{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery := tt.clause.SqlString()
			gotValues := tt.clause.Values()
			if !isStrEq(gotQuery, tt.wantQuery) {
				t.Errorf("SqlString() = %v, want \n%v", gotQuery, removeNewLine(tt.wantQuery))
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestDeleteClauses(t *testing.T) {
	tests := []struct {
		name       string
		clause     *deleteClause
		wantQuery  string
		wantValues []interface{}
	}{
		{
			name: "simple",
			clause: DeleteFrom(
				Channels,
			).
				Where(Channels.Id.EqInt(1)),
			wantQuery:  `DELETE FROM channels WHERE channels.id = %s`,
			wantValues: []interface{}{1},
		},
		{
			name: "using",
			clause: DeleteFrom(
				Channels,
			).
				Using(Products).
				Where(Products.ChannelId.Eq(Channels.Id), Products.Id.EqInt(1)),
			wantQuery:  `DELETE FROM channels USING products WHERE products.channel_id = channels.id AND products.id = %s`,
			wantValues: []interface{}{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery := tt.clause.SqlString()
			gotValues := tt.clause.Values()
			if !isStrEq(gotQuery, tt.wantQuery) {
				t.Errorf("SqlString() = %v, want \n%v", gotQuery, removeNewLine(tt.wantQuery))
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestIfBranch(t *testing.T) {
	col1 := Column("if")
	col2 := Column("elseif")
	col3 := Column("else")
	br := If(true, col1).ElseIf(false, col2).Else(col3)
	if br.SqlString() != col1.SqlString() {
		t.Fatalf("if stmt error, got:%v, expected:%v", br.SqlString(), col1.SqlString())
		return
	}

	br = If(false, col1).ElseIf(true, col2).Else(col3)
	if br.SqlString() != col2.SqlString() {
		t.Fatalf("if stmt error, got:%v, expected:%v", br.SqlString(), col2.SqlString())
		return
	}

	br = If(false, col1).ElseIf(false, col2).Else(col3)
	if br.SqlString() != col3.SqlString() {
		t.Fatalf("if stmt error, got:%v, expected:%v", br.SqlString(), col3.SqlString())
		return
	}

	br = If(true, col1).ElseIf(true, col2).Else(col3)
	if br.SqlString() != col1.SqlString() {
		t.Fatalf("if stmt error, got:%v, expected:%v", br.SqlString(), col1.SqlString())
		return
	}
}
