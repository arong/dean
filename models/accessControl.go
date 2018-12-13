package models

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
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

// LoginInfo store the login info
type LoginRequest struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type accessControl struct {
	loginMap        map[string]*LoginInfo
	tokenMap        map[string]*LoginInfo
	store           *badger.DB
	defaultPassword string // default password for student
}

type resetPassReq struct {
	Password string
}

func init() {
	Ac.tokenMap = make(map[string]*LoginInfo)
}

// SetStore init handler
func (ac *accessControl) SetStore(db *badger.DB) {
	ac.store = db
}

// LoadToken load all authorised user info
func (ac *accessControl) LoadToken() {
	err := ac.store.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			loginInfo := LoginInfo{}
			err := item.Value(func(v []byte) error {
				err := json.Unmarshal(v, &loginInfo)
				//fmt.Printf("key=%s, value=%s\n", k, v)
				return err
			})
			if err != nil {
				logs.Warn("[] invalid login info")
				return err
			}
			ac.tokenMap[string(k)] = &loginInfo
		}
		return nil
	})

	if err != nil {
		logs.Error("[accessControl::LoadToken] load token error", "err", err)
	}
}

// IsValidPassword check too see if password is valid
func (ac *accessControl) IsValidPassword(p string) error {
	if len(p) < 8 {
		return ErrTooShort
	}

	if len(p) >= 256 {
		return ErrTooLong
	}

	strength := zxcvbn.PasswordStrength(p, nil)
	if strength.Score < 2 {
		logs.Debug("[accessControl::IsValidPassword] password too weak, score", strength.Score)
		return ErrTooWeak
	}

	return nil
}

// EncryptPassword encrypt password
func (ac *accessControl) EncryptPassword(p string) (string, error) {
	encrypted := ""
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		logs.Warn("[accessControl::EncryptPassword] encrypt failed", err)
		return encrypted, err
	}
	encrypted = string(hash)
	logs.Debug("[accessControl::EncryptPassword] hash", encrypted)
	return encrypted, nil
}

type LoginInfo struct {
	UserType     int
	ID           int64
	LoginName    string
	Password     string
	ExpireTime   time.Time // expire time of the token
	CurrentToken string
}

// Login: authorise user and issue token
func (ac *accessControl) Login(req *LoginRequest) (string, error) {
	token := ""

	l, ok := ac.loginMap[req.LoginName]
	if !ok {
		logs.Debug("[accessControl::Login] User not found", req.LoginName)
		return token, errNotExist
	}

	if req.Password != l.Password {
		logs.Info("[accessControl::Login] password not match")
		return token, errPermission
	}

	if l.CurrentToken != "" {
		logs.Debug("[accessControl::Login] remove token", l.CurrentToken)
		delete(ac.tokenMap, l.CurrentToken)
		ac.removeToken(l.CurrentToken)
	}

	token = uuid.New().String()
	l.CurrentToken = token
	ac.tokenMap[token] = l

	// 保存token
	ac.storeToken(l)

	return token, nil
}

// UpdatePassword will update current user and request a new login
func (ac *accessControl) UpdatePassword(req *LoginInfo) error {
	l, ok := ac.loginMap[req.LoginName]
	if !ok {
		logs.Debug("[accessControl::UpdatePassword] user not found", req.LoginName)
		return errNotExist
	}

	l.Password = req.Password

	if l.CurrentToken != "" {
		delete(ac.tokenMap, l.CurrentToken)
		ac.removeToken(l.CurrentToken)
		l.CurrentToken = ""
	}
	return nil
}

// ResetAllStudentPassword reset all students' password to default value
func (ac *accessControl) ResetAllStudentPassword(req *resetPassReq) error {
	if req.Password == ac.defaultPassword {
		return errNotExist
	}

	ac.defaultPassword = req.Password
	return nil
}

func (ac *accessControl) storeToken(l *LoginInfo) {
	buffer, err := json.Marshal(l)
	if err != nil {
		logs.Error("[accessControl::storeToken] fatal error", err)
	}

	err = ac.store.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(l.CurrentToken), buffer)
		return err
	})

	if err != nil {
		logs.Error("[accessControl::storeToken] failed to set value: %v", err)
		return
	}
	logs.Info("[accessControl::storeToken] store token success")
}

func (ac *accessControl) removeToken(token string) {
	err := ac.store.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(token))
		return err
	})

	if err != nil {
		logs.Error("[accessControl::removeToken] failed to set value: %v", err)
		return
	}
}

// VerifyToken check to see if the token is valid
func (ac *accessControl) VerifyToken(token string) (LoginInfo, bool) {
	l, ok := ac.tokenMap[token]
	if !ok {
		return LoginInfo{}, false
	}
	return *l, ok
}

// Logout: logout current user from system
func (ac *accessControl) Logout(token string) error {
	delete(ac.tokenMap, token)
	return nil
}
