module github.com/morgine/moon

go 1.15

require (
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.4.2
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gorm.io/gorm v1.20.8
	github.com/morgine/pkg v0.0.0-20201215094710-dd28233bfdf4
)

replace github.com/morgine/pkg => ../pkg
