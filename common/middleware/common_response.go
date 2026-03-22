package middleware

import (
	"bytes"
	"dropshipbe/common/response"
	"encoding/json"
	"net/http"
)

// customResponseWriter vẫn giữ vai trò chặn và lưu dữ liệu lại
type customResponseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *customResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

// BuildCommonResponse là Middleware tích hợp package response của bạn
func BuildCommonResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Khởi tạo writer chặn dữ liệu
		cw := &customResponseWriter{
			ResponseWriter: w,
			body:           bytes.NewBuffer(nil),
			statusCode:     http.StatusOK,
		}

		// 2. Chuyển tiếp yêu cầu
		next(cw, r)

		// 3. Lấy dữ liệu thô
		grpcData := cw.body.Bytes()

		// 4. Phân loại và gọi package response của bạn
		if cw.statusCode >= 200 && cw.statusCode < 300 {
			// TRƯỜNG HỢP THÀNH CÔNG
			var data json.RawMessage
			if len(grpcData) > 0 {
				data = json.RawMessage(grpcData) // Giữ nguyên cấu trúc JSON
			}

			response.Success(w, data)
		} else {
			// TRƯỜNG HỢP LỖI
			msg := string(grpcData)
			if msg == "" {
				msg = "Có lỗi xảy ra từ máy chủ"
			}

			// Gọi hàm Error từ package của bạn
			response.Error(w, cw.statusCode, msg)
		}
	}
}
