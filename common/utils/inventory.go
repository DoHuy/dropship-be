package utils

import (
	"context"
	"dropshipbe/common/constant"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

// ==========================================
// 1. LUA SCRIPT VÀ KHỞI TẠO BIẾN TOÀN CỤC
// ==========================================

// Script Lua: Kiểm tra tồn kho và trừ đi một cách nguyên tử (Atomic)
const deductStockScript = `
local stockKey = KEYS[1]
local decrementAmount = tonumber(ARGV[1])

local currentStock = tonumber(redis.call('get', stockKey))
if currentStock == nil then
    return -1
end

if currentStock >= decrementAmount then
    redis.call('decrby', stockKey, decrementAmount)
    return 1
else
    return 0
end
`

// [TỐI ƯU HIỆU SUẤT]: Khai báo Script Object ở cấp độ Package.
// go-zero sẽ tự động dùng lệnh EVALSHA để cache script này lên server Redis,
// tránh việc phải gửi và biên dịch lại chuỗi string mỗi lần có khách mua hàng.
var deductStockLua = redis.NewScript(deductStockScript)

// ==========================================
// 2. CÁC HÀM THAO TÁC VỚI KHO (RUNTIME)
// ==========================================

// DeductInventory thực thi lệnh trừ tồn kho siêu tốc trên Redis.
// Trả về:
//
//	 1 : Thành công
//	 0 : Hết hàng
//	-1 : Lỗi hoặc không tìm thấy sản phẩm trong Cache
func DeductInventory(ctx context.Context, rds *redis.Redis, ProductVariantID uint64, quantity int32) (int, error) {
	key := constant.ProductVariantStock(fmt.Sprintf("%d", ProductVariantID))

	// Gọi hàm ScriptRunCtx từ đối tượng Redis Client của go-zero
	resp, err := rds.ScriptRunCtx(ctx, deductStockLua, []string{key}, quantity)
	if err != nil {
		return -1, fmt.Errorf("lỗi thực thi Lua script trên Redis: %v", err)
	}

	// Ép kiểu kết quả trả về từ Lua
	result, ok := resp.(int64)
	if !ok {
		return -1, fmt.Errorf("không thể ép kiểu kết quả từ Redis, nhận được: %v", resp)
	}

	return int(result), nil
}

// RollbackInventory hoàn trả lại số lượng kho nếu các bước phía sau (Lưu DB, gọi API PayPal) bị đứt gánh.
func RollbackInventory(ctx context.Context, rds *redis.Redis, productVariantID uint64, quantity int32) error {
	key := constant.ProductVariantStock(fmt.Sprintf("%d", productVariantID))

	// Hoàn trả lại số lượng bằng lệnh INCRBY
	_, err := rds.IncrbyCtx(ctx, key, int64(quantity))
	if err != nil {
		return fmt.Errorf("lỗi khi hoàn trả tồn kho cho %d: %v", productVariantID, err)
	}

	return nil
}

// ==========================================
// 3. NẠP TỒN KHO LÚC KHỞI ĐỘNG (PRE-WARMING)
// ==========================================

// LoadAllInventoryToRedis nạp toàn bộ tồn kho từ Database lên Redis lúc khởi động Server.
func LoadAllInventoryToRedis(ctx context.Context, rds *redis.Redis, db *gorm.DB) error {
	type Variant struct {
		ID            uint64
		StockQuantity int
	}

	var variants []Variant

	// Truy vấn các variant đang kinh doanh (is_active = true)
	if err := db.WithContext(ctx).Table("variants").Where("is_active = ?", true).Select("id, stock_quantity").Find(&variants).Error; err != nil {
		return fmt.Errorf("lỗi khi query DB lấy tồn kho: %v", err)
	}

	// [ĐÃ SỬA LỖI]: Bỏ dấu "_, " đi, chỉ nhận lại 1 biến err
	err := rds.PipelinedCtx(ctx, func(pipe redis.Pipeliner) error {
		for _, v := range variants {
			variantID := strconv.FormatUint(v.ID, 10)
			key := constant.ProductVariantStock(variantID)

			// TUYỆT ĐỐI QUAN TRỌNG: Tham số thứ 3 là 0 nghĩa là TTL = 0 (Không bao giờ hết hạn)
			pipe.Set(ctx, key, v.StockQuantity, 0)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("lỗi khi nạp tồn kho qua Pipeline: %v", err)
	}

	return nil
}
