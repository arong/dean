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

const (
	hmacSecret = "iLoveLFLSS"
)

type LoginInfo struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type accessControl struct {
	teacherMap map[string]string // user name -> password
	studentMap map[string]string // user name -> token
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

func (ac *accessControl) GenToken(name string) (string, error) {
	return "",nil
	//curr, err := Um.GetUserByName(name)
	//if err != nil {
	//	logs.Error("bug found")
	//	return "", err
	//}
	//
	//// Create a new token object, specifying signing method and the claims
	//// you would like it to contain.
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"id":     curr.ID,
	//	"type":   "teacher",
	//	"expire": time.Now().Add(40 * time.Minute).Unix(),
	//})
	//
	//// Sign and get the complete encoded token as a string using the secret
	//tokenString, err := token.SignedString([]byte(hmacSecret))
	//if err != nil {
	//	logs.Error("failed to generate access token", err)
	//	return "", err
	//}
	//// flush cache
	//return tokenString, nil
}

func (ac *accessControl) Login(name, password string, uType int) (string, error) {
	token := ""
	encrypted := ""
	ok := false

	if uType == TypeStudent {
		encrypted, ok = ac.studentMap[name]
		if !ok {
			return token, errNotExist
		}
	} else if uType == TypeTeacher {
		encrypted, ok = ac.teacherMap[name]
		if !ok {
			return token, errNotExist
		}
	} else {
		return token, errors.New("invalid request")
	}
	logs.Debug(encrypted, password)
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	if err != nil {
		logs.Info("failed", err)
		return token, ErrPasswordError
	}
	return ac.GenToken(name)
}

func (ac *accessControl) Logout(token string) error {
	return nil
}

func (ac *accessControl) AddUser(name, password string) {
	ac.studentMap[name] = password
}
func (ac *accessControl) AddTeacher(name, password string) {
	ac.teacherMap[name] = password
}
