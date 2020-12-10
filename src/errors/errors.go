package errors

import "fmt"

const Others Code = 0

const (
	UsernameAlreadyRegistered   Code = 6001
	UsernameOrPasswordIncorrect Code = 6002
	GoogleAuthCodeIncorrect     Code = 6003
)

var Texts = map[Code]string{
	Others:                      "其他错误",
	UsernameAlreadyRegistered:   "用户名已注册",
	UsernameOrPasswordIncorrect: "用户名或密码错误",
	GoogleAuthCodeIncorrect:     "谷歌验证码出错",
}

type Code int

func (c Code) Error() string {
	return fmt.Sprintf("code: %d, error: %s", c, Texts[c])
}
