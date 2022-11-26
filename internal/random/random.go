package random

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/sethvargo/go-password/password"
	"regexp"
	"strings"
)

const (
	MinPasswdLength = 6
	num             = `[0-9]{1}`
	a_z             = `[a-z]{1}`
	A_Z             = `[A-Z]{1}`
	symbol          = `[!@#~$%^&*()+|_]{1}`
)

// UsernameGenerator return a random username
func UsernameGenerator(gender int) string {
	return strings.ToLower(randomdata.FirstName(gender))
}

func PasswordGenerator() (string, error) {
	// Generate a password that is 10 characters long with 4 digits, 2 symbols and 4 letters,
	// allowing upper case letters, disallowing repeat characters.
	return password.Generate(10, 4, 2, false, false)
}

func MustPasswordGenerator() string {
	return password.MustGenerate(10, 4, 2, false, false)
}

// VerifyPasswordStrength used to  Verifying Password Strength
func VerifyPasswordStrength(p string) error {
	// verify the length of password
	if len(p) < MinPasswdLength {
		return fmt.Errorf("password len is < %d", MinPasswdLength)
	}

	if b, err := regexp.MatchString(num, p); !b || err != nil {
		return fmt.Errorf("password need num :%v", err)
	}

	if b, err := regexp.MatchString(a_z, p); !b || err != nil {
		return fmt.Errorf("password need a_z :%v", err)
	}

	if b, err := regexp.MatchString(A_Z, p); !b || err != nil {
		return fmt.Errorf("password need A_Z :%v", err)
	}

	/*
		if b, err := regexp.MatchString(symbol, p); !b || err != nil {
			return fmt.Errorf("password need symbol :%v", err)
		}
	*/

	return nil
}
