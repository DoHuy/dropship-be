package logic

import (
	"context"
	"encoding/json"
	"time"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"
	model "dropshipbe/model/schema"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
)

type GetFrequentlyBoughtProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFrequentlyBoughtProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFrequentlyBoughtProductsLogic {
	return &GetFrequentlyBoughtProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFrequentlyBoughtProductsLogic) GetFrequentlyBoughtProducts(in *dropshipbe.GetFrequentlyBoughtProductsRequest) (*dropshipbe.ProductListResponse, error) {

	frequentlyBoughts, err := l.svcCtx.EcommerceRepo.GetFrequentlyBoughtProducts(l.ctx, in)
	if err != nil {
		l.Logger.Errorf("Lỗi khi lấy danh sách mua kèm cho ProductID %d: %v", in.ProductId, err)
		return nil, err
	}

	// 2. Chuyển đổi dữ liệu từ Model sang Protobuf Message
	var productItems []*dropshipbe.Product
	for _, p := range frequentlyBoughts {
		productItems = append(productItems, &dropshipbe.Product{
			Id:          p.ID,
			CountryCode: p.Product.Country.Code,
			Name:        p.Product.Name,
			Slug:        p.Product.Slug,
			WowDelay:    "",
			Metadata: map[string]string{
				"metadata": p.Product.Metadata.String(),
			},
			Description: p.Product.Description,
			Rating:      float32(p.Product.Rating),
			ReviewCount: int32(p.Product.ReviewCount),
			IsFeatured:  p.Product.IsFeatured,
			IsTrending:  p.Product.IsTrending,
			IsNew:       p.Product.IsNew,
			Price:       float32(p.Product.Price),

			Status: p.Product.Status,
			// Categories, Images, PriceTiers có thể được thêm vào sau nếu cần
			Categories:        l.convertCategories(p.Product.Categories),
			Galleries:         l.convertGaleries(p.Product.Images),
			ProductPriceTiers: l.convertPriceTiers(p.Product.PriceTiers),
			DescriptionImages: l.convertGaleries(p.Product.Images),
			Options:           l.convertOptions(p.Product.Options),
			Variants:          l.convertVariants(p.Product.Variants),
			MetaTitle:         p.Product.MetaTitle,
			MetaDescription:   p.Product.MetaDescription,
			Vendor:            p.Product.Vendor,
			ProductType:       p.Product.ProductType,
			Badge:             *p.Product.Badge,
			SaleLabel:         *p.Product.SaleLabel,
			SaleTag:           *p.Product.SaleTag,
			FlashSaleEndTime:  p.Product.FlashSaleEndTime.Format(time.RFC3339),
			Sold:              int32(p.Product.Sold),
			Tags:              l.convertTags(p.Product.Tags),
			QuantityEnabled:   p.Product.QuantityEnabled,
			QuickShop:         p.Product.QuickShop,
			CreatedAt:         p.CreatedAt.Format(time.RFC3339),
		})
	}

	// 3. Trả về kết quả cho Client
	return &dropshipbe.ProductListResponse{
		Products: productItems,
	}, nil
}

// --- Products ---
func (l *GetFrequentlyBoughtProductsLogic) convertCategories(categories []model.Category) []*dropshipbe.Category {
	var categoryItems []*dropshipbe.Category
	for _, c := range categories {
		categoryItems = append(categoryItems, &dropshipbe.Category{
			Id:          c.ID,
			Name:        c.Name,
			CountryCode: c.CountryCode,
			Slug:        c.Slug,
			Description: c.Description,
			ImageUrl:    "",
			IsActive:    *c.IsActive,
			ButtonText:  "",
			Alt:         "",
		})
	}
	return categoryItems
}

func (l *GetFrequentlyBoughtProductsLogic) convertGaleries(images []model.ProductImage) []*dropshipbe.Gallery {
	var imageItems []*dropshipbe.Gallery
	expirationDuration := time.Duration(l.svcCtx.Config.R2.LinkExpiration) * time.Minute
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, i := range images {

		// Tạo presigned Image URL
		presignedReq, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
			Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
			Key:    aws.String(i.ImageURL), // luu image key vào trường ImageURL của model
		}, s3.WithPresignExpires(expirationDuration))

		if err != nil {
			logx.Errorf("Lỗi khi tạo presigned URL cho image %s: %v", i.ImageURL, err)
			continue // Bỏ qua ảnh này nếu có lỗi
		}

		presignedVideoReq, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
			Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
			Key:    aws.String(i.VideoURL), // luu video key vào trường VideoURL của model
		}, s3.WithPresignExpires(expirationDuration))

		if err != nil {
			logx.Errorf("Lỗi khi tạo presigned URL cho video %s: %v", i.VideoURL, err)
			continue // Bỏ qua video này nếu có lỗi
		}

		imageItems = append(imageItems, &dropshipbe.Gallery{
			Id:       i.ID,
			ImageUrl: presignedReq.URL,
			VideoUrl: presignedVideoReq.URL,
			AltText:  i.AltText,
			Position: int32(i.Position),
		})

	}
	return imageItems
}

func (l *GetFrequentlyBoughtProductsLogic) convertPriceTiers(priceTiers []model.ProductPriceTier) []*dropshipbe.PriceTier {
	var priceTierItems []*dropshipbe.PriceTier
	for _, pt := range priceTiers {
		priceTierItems = append(priceTierItems, &dropshipbe.PriceTier{
			Id:        pt.ID,
			ProductId: pt.ProductID,
			Price:     float32(pt.Price),
			Savings:   pt.SavingsText,
			Qty:       int32(pt.Qty),
			Tag:       pt.Tag,
			TagClass:  pt.TagClass,
			CreatedAt: pt.CreatedAt.Format(time.RFC3339),
		})
	}
	return priceTierItems
}

func (l *GetFrequentlyBoughtProductsLogic) convertOptions(options []model.Option) []*dropshipbe.Option {
	var optionItems []*dropshipbe.Option
	for _, o := range options {
		var optionValueItems []*dropshipbe.OptionValue
		for _, ov := range o.OptionValues {
			optionValueItems = append(optionValueItems, &dropshipbe.OptionValue{
				Id:        ov.ID,
				Value:     ov.Value,
				ColorCode: ov.ColorCode,
				OptionId:  ov.OptionID,
			})
		}
		optionItems = append(optionItems, &dropshipbe.Option{
			Id:           o.ID,
			Name:         o.Name,
			Code:         o.Code,
			OptionValues: optionValueItems,
		})
	}
	return optionItems
}

func (l *GetFrequentlyBoughtProductsLogic) convertVariants(variants []model.Variant) []*dropshipbe.Variant {
	var variantItems []*dropshipbe.Variant

	expirationDuration := time.Duration(l.svcCtx.Config.R2.LinkExpiration) * time.Minute
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	for _, v := range variants {
		var variantOptionValueItems []*dropshipbe.VariantOption
		for _, ov := range v.OptionValues {
			variantOptionValueItems = append(variantOptionValueItems, &dropshipbe.VariantOption{
				OptionId:      ov.OptionID,
				OptionCode:    ov.ColorCode,
				OptionValueId: ov.ID,
				OptionValue:   ov.Value,
			})
		}

		presignedImage, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
			Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
			Key:    aws.String(v.ImageURL),
		}, s3.WithPresignExpires(expirationDuration))

		if err != nil {
			l.Logger.Errorf("Lỗi khi tạo presigned URL cho image variant %s: %v", v.ImageURL, err)
			continue // Bỏ qua video này nếu có lỗi
		}
		variantItems = append(variantItems, &dropshipbe.Variant{
			Id:             v.ID,
			Sku:            v.Sku,
			ProductId:      v.ProductID,
			ImageUrl:       presignedImage.URL,
			Price:          float32(v.Price),
			Barcode:        v.Barcode,
			CompareAtPrice: float32(v.CompareAtPrice),
			CostPrice:      float32(v.CostPrice),
			StockQuantity:  int32(v.StockQuantity),
			Options:        variantOptionValueItems,
			IsActive:       *v.IsActive,
			CreatedAt:      v.CreatedAt.Format(time.RFC3339),
		})
	}
	return variantItems
}

func (l *GetFrequentlyBoughtProductsLogic) convertTags(jsonData datatypes.JSON) []string {
	var tags []string
	err := json.Unmarshal(jsonData, &tags)
	if err != nil {
		logx.Errorf("Lỗi khi chuyển đổi tags: %v", err)
		return []string{}
	}
	return tags
}
