# ==========================================
# STAGE 1: BUILDER
# ==========================================
FROM golang:1.22-alpine AS builder

# Cài đặt git và timezone data
RUN apk add --no-cache git tzdata

WORKDIR /app

# Copy file quản lý thư viện và tải xuống (giúp cache layer này, build lại rất nhanh)
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn vào container
COPY . .

# Khai báo tham số để biết đang build service nào (vd: mq/mq.go)
ARG BUILD_TARGET
# Tên file thực thi sau khi build (mặc định là app_service)
ARG BINARY_NAME=app_service

# Build ra file binary (CGO_ENABLED=0 giúp chạy mượt trên Alpine)
# Tối ưu dung lượng bằng cờ -ldflags="-s -w" (bỏ qua debug info)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/${BINARY_NAME} ${BUILD_TARGET}

# ==========================================
# STAGE 2: PRODUCTION RUNNER
# ==========================================
# Sử dụng Alpine làm base image cho production (cực nhẹ, chỉ ~5MB)
FROM alpine:latest

# BẮT BUỘC: Cài chứng chỉ SSL để gọi API Telegram/AWS SES không bị lỗi bảo mật
RUN apk --no-cache add ca-certificates tzdata

# Đặt múi giờ mặc định cho log hệ thống
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

# Copy file binary từ Stage 1 sang
COPY --from=builder /app/app_service .

# Copy toàn bộ thư mục cấu hình sang container
COPY etc/ ./etc/

# Command mặc định (Sẽ bị ghi đè bởi docker-compose)
CMD ["./app_service"]