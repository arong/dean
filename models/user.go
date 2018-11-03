package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"sort"
	"sync"
)

var Um userManager

type userManager struct {
	nameMap map[string]*User
	idMap   map[UserID]*User
	mutex   sync.Mutex
}

func (um *userManager) Init(userMap map[UserID]*User) {
	um.nameMap = make(map[string]*User)
	if userMap != nil {
		um.idMap = userMap
		for _, v := range userMap {
			if _, ok := um.nameMap[v.LoginName]; ok {
				logs.Warn("duplicated user found")
			}
			um.nameMap[v.LoginName] = v
		}
	} else {
		um.idMap = make(map[UserID]*User)
	}
}

type User struct {
	StudentID UserID `json:"student_id"`
	Filter
	profile
	LoginInfo
	RegisterID string `json:"register_id"` // 学号
}

type userList []*User

func (cl userList) Len() int {
	return len(cl)
}

func (cl userList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

func (cl userList) Less(i, j int) bool {
	return cl[i].StudentID < cl[j].StudentID
}

type profile struct {
	Age       int    `json:"age"`
	Gender    int    `json:"gender"`
	RealName  string `json:"real_name"`
	Mobile    string `json:"mobile"`
	Address   string `json:"address"`
	Birthday  string `json:"birthday"`
}

type UserID int64

func (um *userManager) AddUser(u *User) (UserID, error) {
	if len(u.LoginName) == 0 {
		return 0, errors.New("invalid name")
	}

	if _, ok := um.nameMap[u.LoginName]; ok {
		return 0, errExist
	}

	if u.Gender < eGenderMale || u.Gender > eGenderUnknown {
		return 0, errors.New("invalid gender")
	}

	err := Ac.IsValidPassword(u.Password)
	if err != nil {
		logs.Debug("invalid password")
		return 0, err
	}

	u.Password, err = Ac.EncryptPassword(u.Password)
	if err != nil {
		logs.Warn("[] encrypt failed", err)
		return 0, err
	}

	err = Ma.InsertUser(u)
	if err != nil {
		logs.Info("add user failed", err)
		return 0, err
	}

	um.nameMap[u.LoginName] = u
	um.idMap[u.StudentID] = u

	Ac.AddUser(u.LoginName, u.Password)
	return u.StudentID, nil
}

func (um *userManager) DelUser(uid UserID) error {
	curr, ok := um.idMap[uid]
	if !ok {
		return errNotExist
	}

	err := Ma.DeleteUser(uid)
	if err != nil {
		logs.Warn("[userManager::DelUser] failed", err)
		return err
	}
	delete(um.nameMap, curr.LoginName)
	delete(um.idMap, uid)
	return nil
}

func (um *userManager) ModUser(u *User) error {
	if u.StudentID == 0 {
		return errors.New("invalid user id")
	}

	curr, ok := um.idMap[u.StudentID]
	if !ok {
		return errNotExist
	}

	if curr.LoginName != u.LoginName {
		curr.LoginName = u.LoginName
	}

	if u.RegisterID != "" && curr.RegisterID != u.RegisterID {
		curr.RegisterID = u.RegisterID
	}
	return nil
}

func (um *userManager) GetUser(uid UserID) (*User, error) {
	if val, ok := um.idMap[uid]; ok {
		return val, nil
	}
	return nil, errors.New("User not exists")
}

func (um *userManager) GetUserByName(name string) (*User, error) {
	val, ok := um.nameMap[name]
	if !ok {
		return nil, errNotExist
	}
	return val, nil
}

func (um *userManager) GetAllUsers(f *Filter) userList {
	ret := userList{}
	for _, v := range um.idMap {
		if f.Grade != 0 && f.Grade != v.Grade {
			continue
		}

		if f.Index != 0 && f.Index != v.Index {
			continue
		}
		ret = append(ret, v)
	}
	sort.Sort(ret)
	return ret
}
