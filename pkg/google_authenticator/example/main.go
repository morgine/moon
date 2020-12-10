package main

import (
	"github.com/morgine/moon/pkg/google_authenticator"
	"log"
	"net/http"
	"strconv"
)

func main() {
	userSecret := "unique user secret" // 每个用户的唯一密钥，可随机生成一个，并绑定到用户账号上，用户通过验证，则表示验证器绑定成功

	// 初始化客户端
	client := google_authenticator.NewClient(google_authenticator.Config{
		QRCodeURIGetter: nil,
		ValidRange:      0,
	})
	// 获得二维码
	http.HandleFunc("/qrcode", func(writer http.ResponseWriter, request *http.Request) {
		user := request.URL.Query().Get("user")
		_, _ = writer.Write([]byte(client.GetQRCodeURI(userSecret, user, 200, 200)))
	})
	// 验证 code
	http.HandleFunc("/check", func(writer http.ResponseWriter, request *http.Request) {
		code := request.URL.Query().Get("code")
		ok, err := client.Verify(userSecret, code)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			_, _ = writer.Write([]byte(strconv.FormatBool(ok)))
		}
	})
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		log.Fatal(err)
	}
}
