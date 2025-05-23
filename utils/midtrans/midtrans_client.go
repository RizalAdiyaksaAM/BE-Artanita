package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"tugas-akhir/config"
)

// Client struct untuk menyimpan ServerKey dan ClientKey dari konfigurasi Midtrans
type Client struct {
	ServerKey string
	ClientKey string
}

// NewClient membuat instansi baru untuk client Midtrans dengan konfigurasi yang diberikan
func NewClient(config config.MidtransConfig) *Client {
	return &Client{
		ServerKey: config.ServerKey,
		ClientKey: config.ClientKey,
	}
}

// VerifyNotificationSignature memverifikasi signature dari webhook Midtrans
func (c *Client) VerifyNotificationSignature(notification Notification) bool {
	// Generate signature untuk memverifikasi webhook
	expectedSignature := c.GenerateSignature(notification)

	// Bandingkan signature yang dikirim oleh Midtrans dengan yang dihasilkan secara lokal
	return notification.SignatureKey == expectedSignature
}

// Perbaikan pada utils/midtrans/client.go
// GenerateSignature menghasilkan signature untuk memverifikasi webhook
func (c *Client) GenerateSignature(notification Notification) string {
	// Midtrans menggunakan format berikut untuk signature:
	// SHA512(order_id + status_code + gross_amount + ServerKey)
	data := notification.OrderID + notification.StatusCode + notification.GrossAmount + c.ServerKey

	// Buat SHA512 hash
	h := sha512.New()
	h.Write([]byte(data))

	// Return signature dalam bentuk hexadecimal
	return hex.EncodeToString(h.Sum(nil))
}
