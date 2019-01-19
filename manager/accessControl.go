package manager

import (
	"encoding/json"
	"time"

	"github.com/juju/ratelimit"

	"github.com/arong/dean/base"

	"github.com/astaxie/beego/logs"
	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var Ac accessControl

var (
	//ErrTooShort      = errors.New("password too short")
	ErrTooLong = errors.New("password too long")
	//ErrTooWeak       = errors.New("password too weak")
	//ErrPasswordError = errors.New("password error")
)

type LoginKey struct {
	UserType  int    `json:"type"`
	LoginName string `json:"login_name"`
}

// LoginInfo store the login info
type LoginRequest struct {
	LoginKey
	Password string `json:"password"`
}

func (l LoginRequest) Check() error {
	if l.LoginName == "" ||
		l.Password == "" ||
		(l.UserType != base.AccountTypeStudent && l.UserType != base.AccountTypeTeacher) {
		return errors.New("invalid parameter")
	}
	return nil
}

type accessControl struct {
	loginMap             map[LoginKey]*LoginInfo
	tokenMap             map[string]*LoginInfo
	store                *badger.DB
	allowDefaultPassword bool
	defaultPassword      string // default password for student
	blackList            map[string]bool
}

type ResetPassReq struct {
	Password string
}

func init() {
	Ac.tokenMap = make(map[string]*LoginInfo)
	Ac.blackList = make(map[string]bool)
}

// SetStore init handler
func (ac *accessControl) SetStore(db *badger.DB) {
	ac.store = db
}

// LoadToken load all authorised user info
func (ac *accessControl) LoadToken() {
	err := ac.store.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			loginInfo := LoginInfo{}
			err := item.Value(func(v []byte) error {
				err := json.Unmarshal(v, &loginInfo)
				return err
			})
			if err != nil {
				logs.Warn("[accessControl::LoadToken] invalid login info", "key", string(k))
				continue
			}
			ac.tokenMap[string(k)] = &loginInfo
		}
		return nil
	})

	if err != nil {
		logs.Error("[accessControl::LoadToken] load token error", "err", err)
	}

	_ = ac.store.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("AllowDefaultPassword"))
		if err != nil {
			return err
		}

		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		ac.allowDefaultPassword = string(valCopy) == "true"
		return nil
	})

	_ = ac.store.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("DefaultPassword"))
		if err != nil {
			return err
		}

		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		ac.defaultPassword = string(valCopy)
		return nil
	})
}

type LoginInfo struct {
	UserType     int
	ID           int64
	LoginName    string
	Password     string
	CurrentToken string
	ExpireTime   time.Time         // expire time of the token
	Bucket       *ratelimit.Bucket // maximum try time
}

// Login: authorise user and issue token
func (ac *accessControl) Login(req *LoginRequest) (string, error) {
	token := ""

	// check black list
	{
		_, ok := ac.blackList[req.LoginName]
		if ok {
			logs.Warn("[accessControl::Login] maybe attack", "login name", req.LoginName)
			return "", errPermission
		}
	}

	l, ok := ac.loginMap[req.LoginKey]
	if !ok {
		if req.UserType != base.AccountTypeStudent {
			logs.Debug("[accessControl::Login] User not found", req.LoginName)
			return token, errNotExist
		}
		student, err := Um.GetStudentByRegisterNumber(req.LoginName)
		if err != nil {
			logs.Debug("[accessControl::Login] student not found", req.LoginName)
			return "", err
		}
		l = &LoginInfo{
			UserType:  base.AccountTypeStudent,
			ID:        student.StudentID,
			LoginName: student.RegisterID,
			Password:  ac.defaultPassword,
		}
		err = Ma.InsertPassword(l)
		if err != nil {
			logs.Debug("[accessControl::Login] InsertPassword failed")
			return "", nil
		}
		ac.loginMap[req.LoginKey] = l
	} else if req.Password != l.Password {
		logs.Info("[accessControl::Login] password not match")
		if l.Bucket == nil {
			l.Bucket = ratelimit.NewBucket(time.Hour, 10)
		}
		if l.Bucket.TakeAvailable(1) <= 0 {
			ac.blackList[l.LoginName] = true
		}
		return token, errPermission
	}

	// remove bucket after success login
	l.Bucket = nil

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

type UpdateRequest struct {
	LoginKey
	Password string
}

func (ac *accessControl) Update(req *UpdateRequest) error {
	l, ok := ac.loginMap[req.LoginKey]
	if !ok {
		logs.Debug("[accessControl::Update] user not exist")
		return errNotExist
	}

	if l.Password == req.Password {
		logs.Info("[accessControl::Update] nothing to do")
		return nil
	}

	err := Ma.UpdatePassword(l.ID, req.Password)
	if err != nil {
		logs.Warn("[accessControl::Update] UpdatePassword failed", "err", err)
		return err
	}

	l.Password = req.Password

	logs.Info("[accessControl::Update] update password success", "loginName", req.LoginName, "password", req.Password)
	return nil
}

// ResetAllStudentPassword reset all students' password to default value
func (ac *accessControl) ResetAllStudentPassword(req *ResetPassReq) error {
	ac.defaultPassword = req.Password
	// drop all students's password in db
	err := Ma.ResetAllPassword(ac.defaultPassword)
	if err != nil {
		logs.Warn("[ResetAllStudentPassword] DropAllPassword failed", "err", err)
		return err
	}

	// flush password
	for _, v := range ac.loginMap {
		if v.UserType == base.AccountTypeStudent {
			v.Password = req.Password
		}
	}

	// set flag
	err = ac.store.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("AllowDefaultPassword"), []byte("true"))
	})

	if err != nil {
		logs.Warn("[ResetAllStudentPassword] set cache failed", "err", err)
	}

	// set password
	err = ac.store.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("DefaultPassword"), []byte(ac.defaultPassword))
	})

	if err != nil {
		logs.Warn("[ResetAllStudentPassword] set cache failed", "err", err)
	}

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
