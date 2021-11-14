package code

import (
	"fmt"
	"strings"

	"github.com/anaskhan96/go-password-encoder"
)

type PasswdCoder struct {
	PwdOpts *password.Options
}

func (c *PasswdCoder) Gen(code string) string {
	salt, encodedPwd := password.Encode(code, c.PwdOpts)

	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
}

func (c *PasswdCoder) Verify(code, encCode string) error {
	pwdInfo := strings.Split(encCode, "$")
	if password.Verify(code, pwdInfo[2], pwdInfo[3], c.PwdOpts) {
		return nil
	} else {
		return fmt.Errorf("cannot verify password")
	}
}
