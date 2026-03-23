package constant

import "fmt"

const (
	ProjectPrefix = "dropship:"

	PrefixProduct  = ProjectPrefix + "product:"
	PrefixCategory = ProjectPrefix + "category:"
	PrefixBlog     = ProjectPrefix + "blog:"
)

// =========================================================================
// Các hàm Helper để sinh ra Cache Key động (có chứa biến)
// =========================================================================

// ProductListByCountryKey tạo key cho danh sách sản phẩm theo quốc gia
func ProductListByCountryKey(countryCode string) string {
	if countryCode == "" {
		return PrefixProduct + "list:all"
	}
	return fmt.Sprintf("%slist:country:%s", PrefixProduct, countryCode)
}

// ProductDetailKey tạo key cho chi tiết 1 sản phẩm (ví dụ dùng cho GetProductBySlug)
func ProductDetailKey(slug string) string {
	return fmt.Sprintf("%sdetail:slug:%s", PrefixProduct, slug)
}

// CategoryListKey tạo key cho danh sách danh mục
func CategoryListKey() string {
	return PrefixCategory + "list:all"
}
