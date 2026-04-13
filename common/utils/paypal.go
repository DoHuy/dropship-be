package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreatePayPalOrder gọi sang PayPal để tạo một Intent thanh toán
func CreatePayPalOrder(paypalBaseURL, clientID, secret, mode string, amount float64) (string, string, error) {

	// 1. Lấy Access Token
	token, err := getPayPalAccessToken(paypalBaseURL, clientID, secret)
	if err != nil {
		return "", "", fmt.Errorf("lỗi lấy token PayPal: %v", err)
	}

	// 2. Build Payload tạo Order
	payload := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]string{
					"currency_code": "USD",
					"value":         fmt.Sprintf("%.2f", amount),
				},
			},
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	// 3. Gọi API tạo Order
	req, _ := http.NewRequest("POST", paypalBaseURL+"/v2/checkout/orders", bytes.NewBuffer(payloadBytes))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("lỗi tạo PayPal order, HTTP status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	orderID := result["id"].(string)

	approveLink := ""
	links := result["links"].([]interface{})
	for _, l := range links {
		link := l.(map[string]interface{})
		if link["rel"] == "approve" {
			approveLink = link["href"].(string)
			break
		}
	}

	return orderID, approveLink, nil
}

// getPayPalAccessToken lấy mã thông báo để xác thực API
func getPayPalAccessToken(baseURL, clientID, secret string) (string, error) {
	req, _ := http.NewRequest("POST", baseURL+"/v1/oauth2/token", bytes.NewBufferString("grant_type=client_credentials"))

	// Basic Auth
	auth := clientID + ":" + secret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+basicAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("không thể lấy token, HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return result["access_token"].(string), nil
}

// CapturePayPalOrder thực hiện lệnh trừ tiền chính thức sau khi khách đã xác nhận (Approve) trên PayPal
func CapturePayPalOrder(paypalBaseURL, clientID, secret, orderID string) (string, map[string]interface{}, error) {
	// 1. Lấy Access Token
	token, err := getPayPalAccessToken(paypalBaseURL, clientID, secret)
	if err != nil {
		return "", nil, fmt.Errorf("lỗi lấy token PayPal: %v", err)
	}

	// 2. Build Request Capture
	url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", paypalBaseURL, orderID)

	// Payload của Capture API có thể để trống hoặc là chuỗi JSON rỗng "{}"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return "", nil, fmt.Errorf("lỗi tạo request capture: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// 3. Thực thi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("lỗi kết nối đến PayPal: %v", err)
	}
	defer resp.Body.Close()

	// 4. Đọc và parse dữ liệu JSON trả về
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// API Capture thường trả về 201 Created (hoặc 200 OK) nếu thành công
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", result, fmt.Errorf("paypal từ chối giao dịch, HTTP Code: %d", resp.StatusCode)
	}

	// 5. Trích xuất trạng thái giao dịch (COMPLETED, PENDING, DECLINED...)
	status := "UNKNOWN"
	if val, ok := result["status"].(string); ok {
		status = val
	}

	// Trả về Trạng thái, toàn bộ chuỗi JSON gốc (để lưu DB đối soát), và lỗi (nếu có)
	return status, result, nil
}
