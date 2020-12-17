package errors

import "fmt"

const (
	StatusUnknown  Code = -1
	StatusOK       Code = 200
	StatusNotFound Code = 404
)

const (
	UsernameAlreadyRegistered   Code = 6001
	UsernameIncorrectFormat     Code = 6002
	PasswordIncorrectFormat     Code = 6003
	UsernameOrPasswordIncorrect Code = 6100
	GoogleAuthCodeIncorrect     Code = 6200
	UserUnauthorized            Code = 6300
)

var Texts = map[Code]string{
	StatusUnknown:               "其他错误",
	StatusOK:                    "OK",
	StatusNotFound:              "Not Found",
	UsernameAlreadyRegistered:   "用户名已注册",
	UsernameIncorrectFormat:     "用户名格式错误",
	PasswordIncorrectFormat:     "密码格式错误",
	UsernameOrPasswordIncorrect: "用户名或密码错误",
	GoogleAuthCodeIncorrect:     "谷歌验证码出错",
	UserUnauthorized:            "用户未登陆",
}

// Code 错误码，紧用于提示前端，前端需要根据业务需要再详细提示用户
type Code int

func (c Code) Error() string {
	return fmt.Sprintf("code: %d, error: %s", c, Texts[c])
}

// Unwrap 尝试将 err 解包为 Code
func Unwrap(err error) (code Code, ok bool) {
	code, ok = err.(Code)
	return
}
