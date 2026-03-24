package constant

import "fmt"

const (
	ProjectPrefix = "dropship:"

	PrefixProduct  = ProjectPrefix + "product:"
	PrefixCategory = ProjectPrefix + "category:"
	PrefixBlog     = ProjectPrefix + "blog:"
)

func ProductListByCountryKey(countryCode string) string {
	if countryCode == "" {
		return PrefixProduct + "list:all"
	}
	return fmt.Sprintf("%slist:country:%s", PrefixProduct, countryCode)
}

func ProductDetailKey(slug string) string {
	return fmt.Sprintf("%sdetail:slug:%s", PrefixProduct, slug)
}

// CategoryListKey tạo key cho danh sách danh mục
func CategoryListKey(country_code string) string {
	if country_code == "" {
		return PrefixCategory + "list:all"
	}
	return fmt.Sprintf("%slist:country:%s", PrefixCategory, country_code)
}

func BannerItemListKey() string {
	return PrefixProduct + "banner:list:all"
}

func BlogPostBySlugKey(slug string, country_code string) string {
	return fmt.Sprintf("%s:country:%s:detail:slug:%s", PrefixBlog, country_code, slug)
}

func BlogPostListByCountryKey(countryCode string) string {
	if countryCode == "" {
		return PrefixBlog + "list:all"
	}
	return fmt.Sprintf("%slist:country:%s", PrefixBlog, countryCode)
}

func FeaturedProductListKey(countryCode string) string {
	if countryCode == "" {
		return PrefixProduct + "featured:list:all"
	}
	return fmt.Sprintf("%sfeatured:list:country:%s", PrefixProduct, countryCode)
}

func NewProductListKey(countryCode string) string {
	if countryCode == "" {
		return PrefixProduct + "new:list:all"
	}
	return fmt.Sprintf("%snew:list:country:%s", PrefixProduct, countryCode)
}
