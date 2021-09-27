package stmt

var (

	// Services is table services
	Services = newServices("services", "")

	// Channels is table channels
	Channels = newChannels("channels", "")

	// Publishers is table publishers
	Publishers = newPublishers("publishers", "")

	// Products is table products
	Products = newProducts("products", "")

	// ProductsStats is table products_stats
	ProductsStats = newProductsStats("products_stats", "")

	// ProductsQueryStats is table products_query_stats
	ProductsQueryStats = newProductsQueryStats("products_query_stats", "")

	// ProductsQueryStatsTmp is table products_query_stats_tmp
	ProductsQueryStatsTmp = newProductsQueryStatsTmp("products_query_stats_tmp", "")

	// ServicesStats is table services_stats
	ServicesStats = newServicesStats("services_stats", "")

	// DownloadCounts is table download_counts
	DownloadCounts = newDownloadCounts("download_counts", "")
)

type services struct {
	table, alias  string
	Id            Column //Id is column of id
	Uid           Column //Uid is column of uid
	PublisherId   Column //PublisherId is column of publisher_id
	Name          Column //Name is column of name
	Homepage      Column //Homepage is column of homepage
	Description   Column //Description is column of description
	Functionality Column //Functionality is column of functionality
	CreatedAt     Column //CreatedAt is column of created_at
	UpdatedAt     Column //UpdatedAt is column of updated_at

}

func newServices(table, alias string) services {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return services{
		table:         table,
		alias:         alias,
		Id:            Column(prefix + ".id"),
		Uid:           Column(prefix + ".uid"),
		PublisherId:   Column(prefix + ".publisher_id"),
		Name:          Column(prefix + ".name"),
		Homepage:      Column(prefix + ".homepage"),
		Description:   Column(prefix + ".description"),
		Functionality: Column(prefix + ".functionality"),
		CreatedAt:     Column(prefix + ".created_at"),
		UpdatedAt:     Column(prefix + ".updated_at"),
	}
}

func (t services) Alias(n string) services {
	return newServices(t.table, n)
}

func (t services) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

type channels struct {
	table, alias string
	Id           Column //Id is column of id
	Uid          Column //Uid is column of uid
	Name         Column //Name is column of name
	CreatedAt    Column //CreatedAt is column of created_at

}

func newChannels(table, alias string) channels {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return channels{
		table:     table,
		alias:     alias,
		Id:        Column(prefix + ".id"),
		Uid:       Column(prefix + ".uid"),
		Name:      Column(prefix + ".name"),
		CreatedAt: Column(prefix + ".created_at"),
	}
}

func (t channels) Alias(n string) channels {
	return newChannels(t.table, n)
}

func (t channels) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

type publishers struct {
	table, alias    string
	Id              Column //Id is column of id
	Uid             Column //Uid is column of uid
	Name            Column //Name is column of name
	Homepage        Column //Homepage is column of homepage
	Description     Column //Description is column of description
	CreatedAt       Column //CreatedAt is column of created_at
	UpdatedAt       Column //UpdatedAt is column of updated_at
	Email           Column //Email is column of email
	IndustryName    Column //IndustryName is column of industry_name
	FoundedDate     Column //FoundedDate is column of founded_date
	FinancingEvents Column //FinancingEvents is column of financing_events
	NameTokens      Column //NameTokens is column of name_tokens

}

func newPublishers(table, alias string) publishers {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return publishers{
		table:           table,
		alias:           alias,
		Id:              Column(prefix + ".id"),
		Uid:             Column(prefix + ".uid"),
		Name:            Column(prefix + ".name"),
		Homepage:        Column(prefix + ".homepage"),
		Description:     Column(prefix + ".description"),
		CreatedAt:       Column(prefix + ".created_at"),
		UpdatedAt:       Column(prefix + ".updated_at"),
		Email:           Column(prefix + ".email"),
		IndustryName:    Column(prefix + ".industry_name"),
		FoundedDate:     Column(prefix + ".founded_date"),
		FinancingEvents: Column(prefix + ".financing_events"),
		NameTokens:      Column(prefix + ".name_tokens"),
	}
}

func (t publishers) Alias(n string) publishers {
	return newPublishers(t.table, n)
}

func (t publishers) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

type products struct {
	table, alias   string
	Id             Column //Id is column of id
	Uid            Column //Uid is column of uid
	PublisherId    Column //PublisherId is column of publisher_id
	ChannelId      Column //ChannelId is column of channel_id
	Platform       Column //Platform is column of platform
	BundleId       Column //BundleId is column of bundle_id
	Name           Column //Name is column of name
	Text           Column //Text is column of text
	Homepage       Column //Homepage is column of homepage
	Description    Column //Description is column of description
	LogoUrl        Column //LogoUrl is column of logo_url
	ScreenshotUrls Column //ScreenshotUrls is column of screenshot_urls
	Discontinued   Column //Discontinued is column of discontinued
	Extra          Column //Extra is column of extra
	CreatedAt      Column //CreatedAt is column of created_at
	UpdatedAt      Column //UpdatedAt is column of updated_at
	Countries      Column //Countries is column of countries

}

func newProducts(table, alias string) products {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return products{
		table:          table,
		alias:          alias,
		Id:             Column(prefix + ".id"),
		Uid:            Column(prefix + ".uid"),
		PublisherId:    Column(prefix + ".publisher_id"),
		ChannelId:      Column(prefix + ".channel_id"),
		Platform:       Column(prefix + ".platform"),
		BundleId:       Column(prefix + ".bundle_id"),
		Name:           Column(prefix + ".name"),
		Text:           Column(prefix + ".text"),
		Homepage:       Column(prefix + ".homepage"),
		Description:    Column(prefix + ".description"),
		LogoUrl:        Column(prefix + ".logo_url"),
		ScreenshotUrls: Column(prefix + ".screenshot_urls"),
		Discontinued:   Column(prefix + ".discontinued"),
		Extra:          Column(prefix + ".extra"),
		CreatedAt:      Column(prefix + ".created_at"),
		UpdatedAt:      Column(prefix + ".updated_at"),
		Countries:      Column(prefix + ".countries"),
	}
}

func (t products) Alias(n string) products {
	return newProducts(t.table, n)
}

func (t products) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

type productsStats struct {
	table, alias  string
	Id            Column //Id is column of id
	Name          Column //Name is column of name
	LogoUrl       Column //LogoUrl is column of logo_url
	Uid           Column //Uid is column of uid
	Extra         Column //Extra is column of extra
	PublisherName Column //PublisherName is column of publisher_name
	PublisherUid  Column //PublisherUid is column of publisher_uid
	ChannelUid    Column //ChannelUid is column of channel_uid
	Country       Column //Country is column of country
	Tags          Column //Tags is column of tags
	Downloads     Column //Downloads is column of downloads
	Paid          Column //Paid is column of paid
	NameTokens    Column //NameTokens is column of name_tokens

}

func newProductsStats(table, alias string) productsStats {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return productsStats{
		table:         table,
		alias:         alias,
		Id:            Column(prefix + ".id"),
		Name:          Column(prefix + ".name"),
		LogoUrl:       Column(prefix + ".logo_url"),
		Uid:           Column(prefix + ".uid"),
		Extra:         Column(prefix + ".extra"),
		PublisherName: Column(prefix + ".publisher_name"),
		PublisherUid:  Column(prefix + ".publisher_uid"),
		ChannelUid:    Column(prefix + ".channel_uid"),
		Country:       Column(prefix + ".country"),
		Tags:          Column(prefix + ".tags"),
		Downloads:     Column(prefix + ".downloads"),
		Paid:          Column(prefix + ".paid"),
		NameTokens:    Column(prefix + ".name_tokens"),
	}
}

func (t productsStats) Alias(n string) productsStats {
	return newProductsStats(t.table, n)
}

func (t productsStats) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

type productsQueryStats struct {
	table, alias  string
	Id            Column //Id is column of id
	Name          Column //Name is column of name
	BundleId      Column //BundleId is column of bundle_id
	PublisherName Column //PublisherName is column of publisher_name
	DownloadsSum  Column //DownloadsSum is column of downloads_sum
	ServicesCount Column //ServicesCount is column of services_count
	Paid          Column //Paid is column of paid
	ServicesUids  Column //ServicesUids is column of services_uids
	ChannelUid    Column //ChannelUid is column of channel_uid
	Tags          Column //Tags is column of tags
	Countries     Column //Countries is column of countries
	NameTokens    Column //NameTokens is column of name_tokens

}

func newProductsQueryStats(table, alias string) productsQueryStats {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return productsQueryStats{
		table:         table,
		alias:         alias,
		Id:            Column(prefix + ".id"),
		Name:          Column(prefix + ".name"),
		BundleId:      Column(prefix + ".bundle_id"),
		PublisherName: Column(prefix + ".publisher_name"),
		DownloadsSum:  Column(prefix + ".downloads_sum"),
		ServicesCount: Column(prefix + ".services_count"),
		Paid:          Column(prefix + ".paid"),
		ServicesUids:  Column(prefix + ".services_uids"),
		ChannelUid:    Column(prefix + ".channel_uid"),
		Tags:          Column(prefix + ".tags"),
		Countries:     Column(prefix + ".countries"),
		NameTokens:    Column(prefix + ".name_tokens"),
	}
}

func (t productsQueryStats) Alias(n string) productsQueryStats {
	return newProductsQueryStats(t.table, n)
}

func (t productsQueryStats) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

type productsQueryStatsTmp struct {
	table, alias  string
	Id            Column //Id is column of id
	Name          Column //Name is column of name
	BundleId      Column //BundleId is column of bundle_id
	PublisherName Column //PublisherName is column of publisher_name
	DownloadsSum  Column //DownloadsSum is column of downloads_sum
	ServicesCount Column //ServicesCount is column of services_count
	Paid          Column //Paid is column of paid
	ServicesUids  Column //ServicesUids is column of services_uids
	ChannelUid    Column //ChannelUid is column of channel_uid
	Tags          Column //Tags is column of tags
	Countries     Column //Countries is column of countries
	NameTokens    Column //NameTokens is column of name_tokens

}

func newProductsQueryStatsTmp(table, alias string) productsQueryStatsTmp {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return productsQueryStatsTmp{
		table:         table,
		alias:         alias,
		Id:            Column(prefix + ".id"),
		Name:          Column(prefix + ".name"),
		BundleId:      Column(prefix + ".bundle_id"),
		PublisherName: Column(prefix + ".publisher_name"),
		DownloadsSum:  Column(prefix + ".downloads_sum"),
		ServicesCount: Column(prefix + ".services_count"),
		Paid:          Column(prefix + ".paid"),
		ServicesUids:  Column(prefix + ".services_uids"),
		ChannelUid:    Column(prefix + ".channel_uid"),
		Tags:          Column(prefix + ".tags"),
		Countries:     Column(prefix + ".countries"),
		NameTokens:    Column(prefix + ".name_tokens"),
	}
}

func (t productsQueryStatsTmp) Alias(n string) productsQueryStatsTmp {
	return newProductsQueryStatsTmp(t.table, n)
}

func (t productsQueryStatsTmp) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

type servicesStats struct {
	table, alias string
	Id           Column //Id is column of id
	Title        Column //Title is column of title
	Uid          Column //Uid is column of uid
	ChannelUid   Column //ChannelUid is column of channel_uid
	Country      Column //Country is column of country
	ProductCount Column //ProductCount is column of product_count
	Downloads    Column //Downloads is column of downloads
	NameTokens   Column //NameTokens is column of name_tokens

}

func newServicesStats(table, alias string) servicesStats {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return servicesStats{
		table:        table,
		alias:        alias,
		Id:           Column(prefix + ".id"),
		Title:        Column(prefix + ".title"),
		Uid:          Column(prefix + ".uid"),
		ChannelUid:   Column(prefix + ".channel_uid"),
		Country:      Column(prefix + ".country"),
		ProductCount: Column(prefix + ".product_count"),
		Downloads:    Column(prefix + ".downloads"),
		NameTokens:   Column(prefix + ".name_tokens"),
	}
}

func (t servicesStats) Alias(n string) servicesStats {
	return newServicesStats(t.table, n)
}

func (t servicesStats) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}

type downloadCounts struct {
	table, alias string
	ChannelId    Column //ChannelId is column of channel_id
	ProductId    Column //ProductId is column of product_id
	CountryCode  Column //CountryCode is column of country_code
	Year         Column //Year is column of year
	Month        Column //Month is column of month
	Count        Column //Count is column of count
	CreatedAt    Column //CreatedAt is column of created_at

}

func newDownloadCounts(table, alias string) downloadCounts {
	prefix := table
	if alias != "" {
		prefix = alias
	}
	return downloadCounts{
		table:       table,
		alias:       alias,
		ChannelId:   Column(prefix + ".channel_id"),
		ProductId:   Column(prefix + ".product_id"),
		CountryCode: Column(prefix + ".country_code"),
		Year:        Column(prefix + ".year"),
		Month:       Column(prefix + ".month"),
		Count:       Column(prefix + ".count"),
		CreatedAt:   Column(prefix + ".created_at"),
	}
}

func (t downloadCounts) Alias(n string) downloadCounts {
	return newDownloadCounts(t.table, n)
}

func (t downloadCounts) SqlString() string {
	if t.alias == "" {
		return t.table
	}
	return t.table + " " + t.alias
}
