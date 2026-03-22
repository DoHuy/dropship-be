package main

import (
	"dropshipbe/common/middleware"
	"dropshipbe/dropshipbe"
	"dropshipbe/dropshipbeclient"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
)

var gatewayConfigFile = flag.String("f", "etc/gateway.yaml", "tệp cấu hình gateway")

func main() {
	flag.Parse()

	var c gateway.GatewayConf
	conf.MustLoad(*gatewayConfigFile, &c)

	gw := gateway.MustNewServer(c)
	defer gw.Stop()

	// 2. Tìm cấu hình gRPC từ Upstreams để tạo Client cho Custom Route
	var rpcConfig zrpc.RpcClientConf
	found := false

	for _, upstream := range c.Upstreams {
		// Kiểm tra nếu Upstream này có cấu hình Grpc
		if upstream.Grpc.Target != "" || len(upstream.Grpc.Endpoints) > 0 || upstream.Grpc.Etcd.Key != "" {
			rpcConfig = *upstream.Grpc
			found = true
			break
		}
	}

	if !found {
		panic("Không tìm thấy cấu hình gRPC Upstream nào trong gateway.yaml")
	}

	// 3. Khởi tạo gRPC Client (Lấy trực tiếp từ config Upstream)
	rpcClient := zrpc.MustNewClient(rpcConfig)
	dropshipSvc := dropshipbeclient.NewDropshipbe(rpcClient)

	// 4. Đăng ký Custom Route cho Form-Data
	gw.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/upload",
				Handler: handleUploadFormData(dropshipSvc),
			},
		},
	)

	// Đăng ký Middleware ở đây:
	gw.Use(middleware.BuildCommonResponse)

	fmt.Printf("Bắt đầu Gateway Server (REST API) tại %s:%d...\n", c.Host, c.Port)
	gw.Start()
}

// handleUploadFormData xử lý việc nhận multipart/form-data từ web và gọi sang gRPC
func handleUploadFormData(svc dropshipbeclient.Dropshipbe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Giới hạn dung lượng toàn bộ request (ví dụ: 32MB)
		// 32 << 20 = 32 * 1024 * 1024 bytes
		const maxMemory = 32 << 20
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			http.Error(w, "Dữ liệu Form không hợp lệ hoặc file quá lớn", http.StatusBadRequest)
			return
		}

		// 2. Lấy danh sách file từ key "files" (khớp với -F "files=@..." trong curl)
		formFiles := r.MultipartForm.File["files"]
		if len(formFiles) == 0 {
			http.Error(w, "Vui lòng chọn ít nhất một file để upload", http.StatusBadRequest)
			return
		}

		var fileDataList []*dropshipbe.FileData

		// 3. Duyệt qua từng file được gửi lên
		for _, header := range formFiles {
			// Mở file từ header
			file, err := header.Open()
			if err != nil {

				return
			}

			// Đọc nội dung file vào mảng byte
			content, err := io.ReadAll(file)
			file.Close() // Đóng file ngay sau khi đọc xong để giải phóng tài nguyên

			if err != nil {
				http.Error(w, "Lỗi khi đọc nội dung file: "+header.Filename, http.StatusInternalServerError)
				return
			}

			// Đóng gói vào struct gRPC
			fileDataList = append(fileDataList, &dropshipbe.FileData{
				Filename: header.Filename,
				Content:  content,
			})
		}

		// 4. Gọi sang gRPC Server
		// Sử dụng r.Context() để nếu client ngắt kết nối (Cancel), gRPC cũng sẽ dừng xử lý (Graceful)
		resp, err := svc.UploadFile(r.Context(), &dropshipbe.UploadFileRequest{
			Files: fileDataList,
		})

		if err != nil {
			// Bạn có thể log err chi tiết ở đây cho server-side
			http.Error(w, "Lỗi hệ thống khi upload: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 5. Trả về kết quả thành công theo format chuẩn {code, msg, data}
		httpx.OkJson(w, resp)
	}
}
