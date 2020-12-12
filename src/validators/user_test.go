package validators_test

import (
	"github.com/morgine/moon/src/errors"
	"github.com/morgine/moon/src/validators"
	"testing"
)

func TestValidUsername(t *testing.T) {
	type testcase struct {
		username string
		err      error
	}
	var testcases = []testcase{
		{"1234567", errors.UsernameMustBeLowercaseLettersAndNumbersLen8To16},
		{"", errors.UsernameMustBeLowercaseLettersAndNumbersLen8To16},
		{"Abcdefghig", errors.UsernameMustBeLowercaseLettersAndNumbersLen8To16},
		{"123456789abcdefg", nil},
		{"123456789abcdefgh", errors.UsernameMustBeLowercaseLettersAndNumbersLen8To16},
		{"12345678", nil},
	}
	for _, tc := range testcases {
		if got, need := validators.User.ValidUsername(tc.username), tc.err; got != need {
			t.Errorf("need: %s, got: %s\n", need, got)
		}
	}
}
