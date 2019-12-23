package htpwd

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Bcrypt generates a hash and returns it. It's that simple :).
func Bcrypt(pw []byte) (hash []byte, cost int, err error) {
	if cost < 4 || cost > 17 {
		return hash, cost, errors.New("Invalid range for bcrypt cost")
	}
	hash, err = bcrypt.GenerateFromPassword([]byte(pw), cost)
	if err != nil {
		return hash, cost, err
	}
	return hash, cost, nil
}
