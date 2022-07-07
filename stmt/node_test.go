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

func TestClauses(t *testing.T) {

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
		clause     *clause
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
				FirstValue(DownloadCounts.Count).Over(PartitionBy(DownloadCounts.CountryCode).OrderBy(DownloadCounts.Year)),
				LastValue(DownloadCounts.Count).Over(PartitionBy(DownloadCounts.CountryCode).OrderBy(DownloadCounts.Year)),
				NthValue(DownloadCounts.Count, 2).Over(PartitionBy(DownloadCounts.CountryCode).OrderBy(DownloadCounts.Year)),
			).
				From(DownloadCounts).
				Where(DownloadCounts.ProductId.EqInt(123)),
			wantQuery: `SELECT 
download_counts.country_code, 
Rank() OVER (PARTITION BY download_counts.country_code ORDER BY download_counts.year) AS year_rank, 
FIRST_VALUE(download_counts.count) OVER (PARTITION BY download_counts.country_code ORDER BY download_counts.year), 
LAST_VALUE(download_counts.count) OVER (PARTITION BY download_counts.country_code ORDER BY download_counts.year), 
NTH_VALUE(download_counts.count, 2) OVER (PARTITION BY download_counts.country_code ORDER BY download_counts.year) 
FROM download_counts 
WHERE download_counts.product_id = %s`,
			wantValues: []interface{}{123},
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
