package validators

import (
	"github.com/morgine/moon/src/errors"
	"regexp"
)

//var User = &user{
//	username:   regexp.MustCompile("^[a-z0-9]{8,16}$"),
//	password: regexp.MustCompile("^[\\w]{8,16}$"),
//}

type User interface {
	ValidUsername(username string) error
	ValidPassword(password string) error
}

func NewUser(username, password *regexp.Regexp) User {
	return &user{
		username: username,
		password: password,
	}
}

type user struct {
	username *regexp.Regexp
	password *regexp.Regexp
}

func (u *user) ValidUsername(username string) error {
	if !u.username.MatchString(username) {
		return errors.UsernameIncorrectFormat
	}
	return nil
}

func (u user) ValidPassword(password string) error {
	if !u.password.MatchString(password) {
		return errors.PasswordIncorrectFormat
	}
	return nil
}
