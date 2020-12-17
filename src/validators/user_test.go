package validators_test

import (
	"github.com/morgine/moon/src/errors"
	"github.com/morgine/moon/src/validators"
	"regexp"
	"testing"
)

func TestUser_ValidUsername(t *testing.T) {
	type testcase struct {
		username string
		err      error
	}
	var testcases = []testcase{
		{"1234567", errors.UsernameIncorrectFormat},
		{"", errors.UsernameIncorrectFormat},
		{"Abcdefghig", errors.UsernameIncorrectFormat},
		{"123456789abcdefg", nil},
		{"123456789abcdefgh", errors.UsernameIncorrectFormat},
		{"12345678", nil},
	}
	user := validators.NewUser(regexp.MustCompile("^[a-z0-9]{8,16}$"), nil)
	for _, tc := range testcases {
		if got, need := user.ValidUsername(tc.username), tc.err; got != need {
			t.Errorf("username: %s, need: %s, got: %s\n", tc.username, need, got)
		}
	}
}

func TestUser_ValidUserPassword(t *testing.T) {
	type testcase struct {
		password string
		err      error
	}
	var testcases = []testcase{
		{"1234567", errors.PasswordIncorrectFormat},
		{"", errors.PasswordIncorrectFormat},
		{"123ABCDEFG", nil},
		{"123456789abcdefg", nil},
		{"123456789abcdefgh", errors.PasswordIncorrectFormat},
		{"12345678", nil},
	}
	user := validators.NewUser(nil, regexp.MustCompile("^[\\w]{8,16}$"))
	for _, tc := range testcases {
		if got, need := user.ValidPassword(tc.password), tc.err; got != need {
			t.Errorf("password: %s, need: %s, got: %s\n", tc.password, need, got)
		}
	}
}
