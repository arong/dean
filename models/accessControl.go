package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/nbutton23/zxcvbn-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var Ac accessControl

var (
	ErrTooShort      = errors.New("password too short")
	ErrTooLong       = errors.New("password too long")
	ErrTooWeak       = errors.New("password too weak")
	ErrPasswordError = errors.New("password error")
)

type accessControl struct {
	secret map[string]string
}

func (ac *accessControl) IsValidPassword(p string) error {
	if len(p) < 8 {
		return ErrTooShort
	}

	if len(p) >= 256 {
		return ErrTooLong
	}

	strength := zxcvbn.PasswordStrength(p, nil)
	if strength.Score < 2 {
		logs.Debug("[] password too weak, score", strength.Score)
		return ErrTooWeak
	}

	return nil
}

func (ac *accessControl) EncryptPassword(p string) (string, error){
	encrypted := ""
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		logs.Warn("[] encrypt failed", err)
		return encrypted, err
	}
	encrypted = string(hash)
	logs.Debug("[] hash", encrypted)
	return encrypted, nil
}

func (ac *accessControl) isMatch(plain, encrypted string) error {
	return nil
}

func (ac *accessControl) Login(name, password string) (string, error) {
	token := ""
	val, ok := ac.secret[name]
	if !ok {
		return token, errNotExist
	}

	if ac.isMatch(password, val) != nil {
		logs.Info("invalid password")
		return token, ErrPasswordError
	}
	return "aronic", nil
}

func (ac *accessControl) Logout(token string) error {
	return nil
}
