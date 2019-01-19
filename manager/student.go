package manager

import (
	"sort"
	"sync"

	"github.com/arong/dean/models"

	"github.com/arong/dean/base"
	"github.com/astaxie/beego/logs"
)

// Um user manager
var Um StudentManager

type StudentManager struct {
	// store current active student
	list models.StudentList
	// student id -> index
	idMap map[int64]int
	// student register number -> index
	uuidMap map[string]int
	mutex   sync.Mutex
	// storage
	store models.StudentStore
}

func (s *StudentManager) save(info models.StudentInfo) {
	info.Status = base.StatusValid
	s.list = append(s.list, info)
	k := len(s.list) - 1
	s.idMap[info.StudentID] = k
	s.uuidMap[info.RegisterID] = k
}

func (s *StudentManager) update(info models.StudentInfo) {
	k, ok := s.idMap[info.StudentID]
	if ok {
		old := s.list[k]
		if old.RegisterID != info.RegisterID {
			delete(s.uuidMap, old.RegisterID)
			s.uuidMap[info.RegisterID] = k
		}
		s.list[k] = info
	}
}

func (s *StudentManager) delete(k int) {
	p := &s.list[k]
	p.Status = base.StatusDeleted
	delete(s.idMap, p.StudentID)
	delete(s.uuidMap, p.RegisterID)
}

func (s *StudentManager) get(id int64) (models.StudentInfo, error) {
	if k, ok := s.idMap[id]; ok {
		return s.list[k], nil
	} else {
		return models.StudentInfo{}, errNotExist
	}
}

func (s *StudentManager) Init(list models.StudentList) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.list = models.StudentList{}
	s.idMap = make(map[int64]int)
	s.uuidMap = make(map[string]int)

	if list != nil {
		for k, v := range list {
			if v.StudentID != base.StatusValid {
				continue
			}
			tmp := v
			s.list = append(s.list, tmp)
			s.idMap[v.StudentID] = k
			s.uuidMap[v.RegisterID] = k
		}
	}
}

func (um *StudentManager) AddUser(u models.StudentInfo) (int64, error) {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	if _, ok := um.uuidMap[u.RegisterID]; ok {
		return 0, errExist
	}

	studentID, err := um.store.SaveStudent(u)
	if err != nil {
		logs.Info("[AddUser]add user failed", err)
		return 0, err
	}

	u.StudentID = studentID

	um.save(u)

	return u.StudentID, nil
}

func (um *StudentManager) UpdateStudent(u models.StudentInfo) error {
	curr, err := um.get(u.StudentID)
	if err != nil {
		return err
	}

	if curr.Equal(u) {
		logs.Info("[] nothing to do")
		return nil
	}

	if u.RegisterID != "" && curr.RegisterID != u.RegisterID {
		curr.RegisterID = u.RegisterID
	}
	return nil
}

func (um *StudentManager) DelUser(uidList []int64) ([]int64, error) {
	failed := []int64{}
	load := []int64{}

	for _, uid := range uidList {
		_, ok := um.idMap[uid]
		if !ok {
			failed = append(failed, uid)
			continue
		}

		load = append(load, uid)
		delete(um.idMap, uid)
	}

	err := um.store.DeleteStudent(load)
	if err != nil {
		failed = append(failed, load...)
		logs.Warn("[StudentManager::DelUser] failed", err)
		return failed, err
	}
	return failed, nil
}

func (um *StudentManager) GetUser(uid int64) (models.StudentInfo, error) {
	if uid <= 0 {
		return models.StudentInfo{}, errInvalidParam
	}
	um.mutex.Lock()
	defer um.mutex.Unlock()

	return um.get(uid)
}

func (um *StudentManager) GetStudentByRegisterNumber(reg string) (models.StudentInfo, error) {
	um.mutex.Lock()
	um.mutex.Unlock()

	s, ok := um.uuidMap[reg]
	if !ok {
		return models.StudentInfo{}, errNotExist
	}
	return um.list[s], nil
}

func (um *StudentManager) GetAllUsers(f models.StudentFilter) base.CommList {
	resp := base.CommList{}
	ret := models.StudentList{}
	total := 0
	start, end := f.GetRange()

	logs.Debug("[GetAllStudent]", "total count", len(um.idMap), "start", start, "end", end)
	sort.Sort(ret)
	resp.List = ret
	resp.Total = total
	return resp
}

func (um *StudentManager) getStudentList(grade int) ([]int64, error) {
	ret := []int64{}
	for _, v := range um.list {
		ret = append(ret, v.StudentID)
	}
	return ret, nil
}

func (um *StudentManager) IsExist(studentID int64) bool {
	um.mutex.Lock()
	um.mutex.Unlock()
	_, ok := um.idMap[studentID]
	return ok
}
