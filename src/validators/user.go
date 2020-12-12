package validators

import (
	"github.com/morgine/moon/src/errors"
	"regexp"
)

var User = &user{
	nameReg:   regexp.MustCompile("[a-z0-9]{8,16}"),
	passwdReg: regexp.MustCompile("[\\w]{8,16}"),
}

type user struct {
	nameReg   *regexp.Regexp
	passwdReg *regexp.Regexp
}

func (u *user) ValidUsername(username string) error {
	if !u.nameReg.MatchString(username) {
		return errors.UsernameMustBeLowercaseLettersAndNumbersLen8To16
	}
	return nil
}

func (u user) ValidUserPassword(password string) error {
	if !u.passwdReg.MatchString(password) {
		return errors.PasswordMustBeLettersAndNumbersLen8To16
	}
	return nil
}
