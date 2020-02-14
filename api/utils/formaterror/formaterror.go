package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "nickname") {
		return errors.New("nickname already taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("email already taken")
	}
	if strings.Contains(err, "title") {
		return errors.New("title already taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect Password")
	}
	return errors.New("incorrect details")
}
