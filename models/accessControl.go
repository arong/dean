package models

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/dgraph-io/badger"
	"github.com/nbutton23/zxcvbn-go"
	"github.com/google/uuid"
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
type LoginInfo struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type accessControl struct {
	loginMap map[string]*loginInfo
	tokenMap map[string]*loginInfo
	store    *badger.DB
}

func init() {
	Ac.tokenMap = make(map[string]*loginInfo)
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
			loginInfo := loginInfo{}
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

type loginInfo struct {
	UserType     int
	ID           int64
	LoginName    string
	Password     string
	ExpireTime   time.Time // expire time of the token
	CurrentToken string
}

// Login: authorise user and issue token
func (ac *accessControl) Login(req *LoginInfo) (string, error) {
	token := ""

	l, ok := ac.loginMap[req.LoginName]
	if !ok {
		logs.Debug("[accessControl::Login] User not found", req.LoginName)
		return token, errNotExist
	}

	err := bcrypt.CompareHashAndPassword([]byte(l.Password), []byte(req.Password))
	if err != nil {
		logs.Info("[accessControl::Login] failed", err)
		return token, ErrPasswordError
	}

	if l.CurrentToken != "" {
		logs.Debug("remove token", l.CurrentToken)
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

func (ac *accessControl) storeToken(l *loginInfo) {
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
func (ac *accessControl) VerifyToken(token string) bool {
	_, ok := ac.tokenMap[token]
	return ok
}

// Logout: logout current user from system
func (ac *accessControl) Logout(token string) error {
	delete(ac.tokenMap, token)
	return nil
}
