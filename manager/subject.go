package manager

import (
	"sort"
	"sync"

	"github.com/arong/dean/models"

	"github.com/arong/dean/base"

	"github.com/astaxie/beego/logs"
)

// Sm is global subject manager
var Sm SubjectManager

// SubjectManager manager all course subject
type SubjectManager struct {
	list    models.SubjectList // actual storage
	idMap   map[int]int        // id -> index
	nameMap map[string]int     // name -> index
	keyMap  map[string]int
	ref     map[int]int // id -> reference count
	mutex   sync.Mutex
	store   models.SubjectStore
}

func (sm *SubjectManager) save(info models.SubjectInfo) {
	info.Status = base.StatusValid
	sm.list = append(sm.list, info)
	k := len(sm.list) - 1
	sm.idMap[info.ID] = k
	sm.keyMap[info.Key] = k
	sm.nameMap[info.Name] = k
}

func (sm *SubjectManager) delete(k int) {
	p := &sm.list[k]
	p.Status = base.StatusDeleted
	delete(sm.idMap, p.ID)
	delete(sm.keyMap, p.Key)
	delete(sm.nameMap, p.Name)
}

func (sm *SubjectManager) update(info models.SubjectInfo) {
	k := sm.idMap[info.ID]
	p := &sm.list[k]
	if p.Key != info.Key {
		delete(sm.keyMap, p.Key)
		sm.keyMap[info.Key] = k
	}

	sm.list[k] = info
}

func (sm *SubjectManager) get(id int) (models.SubjectInfo, error) {
	if val, ok := sm.idMap[id]; ok {
		return sm.list[val], nil
	} else {
		return models.SubjectInfo{}, errNotExist
	}
}

func (sm *SubjectManager) Init(list models.SubjectList) {
	sm.idMap = make(map[int]int)
	sm.keyMap = make(map[string]int)
	sm.nameMap = make(map[string]int)
	sm.ref = make(map[int]int)
	sm.list = models.SubjectList{}

	if list != nil {
		for k, v := range list {
			tmp := v
			sm.list = append(sm.list, tmp)
			sm.idMap[v.ID] = k
			sm.keyMap[v.Key] = k
			sm.nameMap[v.Name] = k
		}
	}
}

func (sm *SubjectManager) IncRef(id int) {
	sm.mutex.Lock()
	if _, ok := sm.idMap[id]; ok {
		sm.ref[id] = sm.ref[id] + 1
	}
	sm.mutex.Unlock()
}

func (sm *SubjectManager) DecRef(id int) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, ok := sm.idMap[id]; ok {
		curr := sm.ref[id]
		if curr <= 0 {
			logs.Warn("[SubjectManager::DecRef] bug found", "id", id, "curr", curr)
			return
		}
		sm.ref[id] = curr - 1
	}
}

// Add add new subject into manager
func (sm *SubjectManager) Add(s models.SubjectInfo) (int, error) {
	if _, ok := sm.nameMap[s.Name]; ok {
		logs.Info("[SubjectManager::Add] name duplicate")
		return 0, errExist
	}

	if _, ok := sm.keyMap[s.Key]; ok {
		logs.Info("[SubjectManager::Add] key duplicate")
		return 0, errExist
	}

	var err error
	s.ID, err = sm.store.SaveSubject(s)
	if err != nil {
		logs.Info("[SubjectManager::Add] db failed", "err", err)
		return 0, err
	}
	// add to current map
	sm.mutex.Lock()
	sm.save(s)
	sm.mutex.Unlock()

	return s.ID, nil
}

// Update update the key of subject
func (sm *SubjectManager) Update(s models.SubjectInfo) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	curr, err := sm.get(s.ID)
	if err != nil {
		logs.Info("[SubjectManager::Update] not found", "id", s.ID)
		return errNotExist
	}

	if curr.Equal(s) {
		logs.Info("[SubjectManager::Update] no change")
		return nil
	}

	if curr.Name != s.Name {
		return errPermission
	}

	if curr.Key != s.Key {
		curr.Key = s.Key
	}

	k := sm.idMap[s.ID]
	// name could not change, for sake of user use
	//if tmp, ok := sm.nameMap[s.Name]; ok && tmp != k {
	//	logs.Info("[SubjectManager::Update] name duplicate")
	//	return errExist
	//}

	if tmp, ok := sm.keyMap[s.Key]; ok && tmp != k {
		logs.Info("[SubjectManager::Update] key duplicate")
		return errExist
	}

	err = sm.store.UpdateSubject(s)
	if err != nil {
		logs.Info("[SubjectManager::Update] db failed", "err", err)
		return err
	}

	// update cache
	sm.update(curr)

	return nil
}

// Delete remove subject of id
func (sm *SubjectManager) Delete(ids []int) ([]int, error) {
	var err error
	failedList := []int{}
	for _, id := range ids {
		k, ok := sm.idMap[id]
		if !ok {
			logs.Info("[SubjectManager::Delete] class not found", "id", id)
			failedList = append(failedList, id)
			continue
		}

		if sm.ref[id] > 0 {
			logs.Info("[SubjectManager::Delete] still in use", "id", id)
			failedList = append(failedList, id)
			continue
		}

		err = sm.store.DeleteSubject(id)
		if err != nil {
			logs.Info("[SubjectManager::Delete] db failed", "err", err)
			return nil, err
		}
		// add to current map
		sm.delete(k)
	}
	return failedList, nil
}

// GetAll return all subject id list in current manager
func (sm *SubjectManager) GetAll() models.SubjectList {
	ret := models.SubjectList{}
	sm.mutex.Lock()
	for _, v := range sm.list {
		if v.Status == base.StatusDeleted {
			continue
		}
		tmp := v
		ret = append(ret, tmp)
	}
	sm.mutex.Unlock()

	sort.Sort(ret)
	return ret
}

// CheckSubjectList check to see if all id in input list exist
func (sm *SubjectManager) CheckSubjectList(list []int) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for _, v := range list {
		if _, ok := sm.idMap[v]; !ok {
			return false
		}
	}
	return true
}

// IsExist check to see if exist
func (sm *SubjectManager) IsExist(id int) bool {
	sm.mutex.Lock()
	_, ok := sm.idMap[id]
	sm.mutex.Unlock()
	return ok
}

func (sm *SubjectManager) getSubjectName(id int) string {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	s, err := sm.get(id)
	if err != nil {
		return ""
	}
	return s.Name
}
