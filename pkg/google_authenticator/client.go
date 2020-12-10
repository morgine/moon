package google_authenticator

import (
	"encoding/base32"
	"github.com/dgryski/dgoogauth"
	"net/url"
	"strconv"
)

type QRCodeURIGetter func(QRCodeContent string, with, height int) string

type Client struct {
	config Config
}

type Config struct {
	QRCodeURIGetter QRCodeURIGetter // 二维码图片及地址生成器，不传该参数则默认通过 https://api.qrserver.com 生成二维码图片地址
	ValidRange      int             // 验证区间，即前面 n 个验证码算作有效验证码，取值范围 0-100，最好不要超过 3
}

func NewClient(c Config) *Client {
	if c.QRCodeURIGetter == nil {
		// 使用第三方二维码生成器
		c.QRCodeURIGetter = func(content string, with, height int) string {
			return "https://api.qrserver.com/v1/create-qr-code/?data=" +
				url.QueryEscape(content) +
				"&size=" + strconv.Itoa(with) + "x" + strconv.Itoa(height) + "&ecc=M"
		}
	}
	return &Client{
		config: c,
	}
}

// 获得二维码内容
func (c *Client) GetQRCodeURI(secret, user string, with, height int) string {
	cfg := &dgoogauth.OTPConfig{
		Secret: base32.StdEncoding.EncodeToString([]byte(secret)),
	}
	return c.config.QRCodeURIGetter(cfg.ProvisionURI(user), with, height)
}

// 验证
func (c *Client) Verify(secret, code string) (bool, error) {
	cfg := &dgoogauth.OTPConfig{
		Secret:     base32.StdEncoding.EncodeToString([]byte(secret)),
		WindowSize: c.config.ValidRange,
	}
	return cfg.Authenticate(code)
}
