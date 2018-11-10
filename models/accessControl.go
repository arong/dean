package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/nbutton23/zxcvbn-go"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var Ac accessControl

var (
	ErrTooShort      = errors.New("password too short")
	ErrTooLong       = errors.New("password too long")
	ErrTooWeak       = errors.New("password too weak")
	ErrPasswordError = errors.New("password error")
)

type LoginInfo struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type accessControl struct {
	loginMap map[string]*loginInfo
	tokenMap map[string]*loginInfo
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

func (ac *accessControl) EncryptPassword(p string) (string, error) {
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

type loginInfo struct {
	UserType  int
	ID        UserID
	LoginName string
	Password  string
}

func (ac *accessControl) Login(req *LoginInfo) (string, error) {
	token := ""

	l, ok := ac.loginMap[req.LoginName]
	if !ok {
		logs.Debug("User not found", req.LoginName)
		return token, errNotExist
	}

	err := bcrypt.CompareHashAndPassword([]byte(l.Password), []byte(req.Password))
	if err != nil {
		logs.Info("failed", err)
		return token, ErrPasswordError
	}

	token = uuid.NewV4().String()
	ac.tokenMap[token] = l

	return token, nil
}

func (ac *accessControl) Logout(token string) error {
	return nil
}

func (ac *accessControl) AddUser(name, password string) {
}
func (ac *accessControl) AddTeacher(name, password string) {
}
