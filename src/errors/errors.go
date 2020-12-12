package errors

import "fmt"

const Others Code = 0

const (
	UsernameAlreadyRegistered                        Code = 6001
	UsernameMustBeLowercaseLettersAndNumbersLen8To16 Code = 6002
	PasswordMustBeLettersAndNumbersLen8To16          Code = 6003
	PasswordLengthMustRange8To16                     Code = 6005
	UsernameOrPasswordIncorrect                      Code = 6100
	GoogleAuthCodeIncorrect                          Code = 6200
)

var Texts = map[Code]string{
	Others:                    "其他错误",
	UsernameAlreadyRegistered: "用户名已注册",
	UsernameMustBeLowercaseLettersAndNumbersLen8To16: "用户名必须是小写字母或数字的组合, 且长度在8-16位之间",
	PasswordMustBeLettersAndNumbersLen8To16:          "密码必须是字母或数字的组合, 且长度在8-16位之间",
	UsernameOrPasswordIncorrect:                      "用户名或密码错误",
	GoogleAuthCodeIncorrect:                          "谷歌验证码出错",
}

type Code int

func (c Code) Error() string {
	return fmt.Sprintf("code: %d, error: %s", c, Texts[c])
}

// Unwrap 尝试将 err 解包为 Code
func Unwrap(err error) (code Code, ok bool) {
	code, ok = err.(Code)
	return
}
