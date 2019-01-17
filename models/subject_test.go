package models

import (
	"log"
	"testing"

	"github.com/pkg/errors"

	"github.com/arong/dean/base"
	"github.com/golang/mock/gomock"
)

var list = SubjectList{
	{Status: base.StatusValid, ID: 81, Name: "语文", Key: "Chinese"},
	{Status: base.StatusValid, ID: 82, Name: "数学", Key: "Math"},
	{Status: base.StatusValid, ID: 83, Name: "英语", Key: "English"},
	{Status: base.StatusValid, ID: 84, Name: "物理", Key: "Physics"},
	{Status: base.StatusValid, ID: 85, Name: "化学", Key: "Chemistry"},
	{Status: base.StatusValid, ID: 86, Name: "生物", Key: "Biology"},
	{Status: base.StatusValid, ID: 87, Name: "政治", Key: "Politics"},
	{Status: base.StatusValid, ID: 88, Name: "历史", Key: "History"},
	{Status: base.StatusValid, ID: 89, Name: "地理", Key: "Geology"},
	{Status: base.StatusValid, ID: 90, Name: "体育", Key: "P.E."},
}

func TestSubjectManager_Init(t *testing.T) {
	sm := SubjectManager{}
	sm.Init(list)
	for _, v := range list {
		if !sm.IsExist(v.ID) {
			t.Fatal("init failed")
		}
	}
}

func TestSubjectManager_IsExist(t *testing.T) {
	sm := SubjectManager{}
	sm.Init(list)

	for _, v := range list {
		if !sm.IsExist(v.ID) {
			t.Fatal("IsExist failed")
		}
	}

	if sm.IsExist(1) {
		t.Fatal("IsExist failed")
	}
}

func TestSubjectManager_IncRef(t *testing.T) {
	sm := SubjectManager{}
	sm.Init(list)

	for _, v := range list {
		sm.IncRef(v.ID)
	}

	load := []int{}
	for _, v := range list {
		load = append(load, v.ID)
	}

	failed, _ := sm.Delete(load)
	if len(failed) != len(list) {
		t.Fatal("IncRef failed")
	}
}

func TestSubjectManager_DecRef(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := NewMockSubjectStore(mockCtrl)
	sm := SubjectManager{store: mockStore}

	sm.Init(list)

	for _, v := range list {
		sm.IncRef(v.ID)
	}

	for _, v := range list {
		sm.DecRef(v.ID)
		sm.DecRef(v.ID)
	}

	for _, v := range sm.ref {
		if v < 0 {
			t.Fatal("bug found")
		}
	}

	load := []int{}
	for _, v := range list {
		load = append(load, v.ID)
	}

	mockStore.EXPECT().DeleteSubject(gomock.Any()).Return(nil).Times(len(list))
	failed, _ := sm.Delete(load)
	if len(failed) != 0 {
		t.Fatal("DecRef failed")
	}
}

func TestSubjectManager_Add(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := NewMockSubjectStore(mockCtrl)
	sm := SubjectManager{store: mockStore}

	sm.Init(list)

	for _, v := range list {
		_, err := sm.Add(v)
		if err == nil {
			t.Fatal("add duplicated success")
		}

		tmp := v
		tmp.Name += "modify"
		_, err = sm.Add(tmp)
		if err == nil {
			t.Fatal("add duplicated key")
		}
	}

	sm.Init(nil)
	for k, v := range list {
		mockStore.EXPECT().SaveSubject(gomock.Any()).Return(k+1, nil)
		_, err := sm.Add(v)
		if err != nil {
			t.Fatal("add failed")
		}
	}

	sm.Init(nil)
	mockStore.EXPECT().SaveSubject(gomock.Any()).Return(0, errors.New("sank your ship"))
	_, err := sm.Add(list[0])
	if err == nil {
		t.Fatal("logic error")
	}
}

func TestSubjectManager_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := NewMockSubjectStore(mockCtrl)
	sm := SubjectManager{store: mockStore}

	sm.Init(list)
	for _, v := range list {
		{
			err := sm.Update(v)
			if err != nil {
				t.Fatal("update failed for no-changes")
			}
		}

		{
			tmp := v
			tmp.Name += "modify"
			err := sm.Update(tmp)
			if err == nil {
				t.Fatal("modify name success")
			}
		}

		{
			tmp := v
			tmp.Key += "modify"
			mockStore.EXPECT().UpdateSubject(gomock.Any()).Return(nil)
			err := sm.Update(tmp)
			if err != nil {
				t.Fatal("modify key failed")
			}
		}
	}

	{
		sm.Init(nil)
		err := sm.Update(list[0])
		if err == nil {
			t.Fatal("update non-existing success")
		}
	}
	{
		sm.Init(list)
		tmp := list[0]
		tmp.Key = list[1].Key
		err := sm.Update(tmp)
		if err == nil {
			t.Fatal("update failure")
		}
	}

	{
		sm.Init(list)
		tmp := list[0]
		tmp.Key += "modified"
		mockStore.EXPECT().UpdateSubject(gomock.Any()).Return(errors.New("sank your ship"))
		err := sm.Update(tmp)
		if err == nil {
			t.Fatal("update failure")
		}
	}
}

func TestSubjectManager_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := NewMockSubjectStore(mockCtrl)
	sm := SubjectManager{store: mockStore}

	sm.Init(list)
	load := []int{}
	for _, v := range list {
		load = append(load, v.ID)
	}

	mockStore.EXPECT().DeleteSubject(gomock.Any()).Return(nil).Times(len(list))
	_, err := sm.Delete(load)
	if err != nil {
		t.Fatal("delete subject success")
	}

	tmp := sm.GetAll()
	if len(tmp) > 0 {
		t.Fatal("delete failed")
	}

	failed, err := sm.Delete(load)
	if len(failed) != len(load) {
		t.Fatal("delete non-existing failed")
	}

	sm.Init(nil)
	mockStore.EXPECT().SaveSubject(gomock.Any()).Return(100, nil)
	id, err := sm.Add(list[0])
	if err != nil || id != 100 {
		t.Fatal("logic failure")
	}
	mockStore.EXPECT().DeleteSubject(gomock.Any()).Return(errors.New("sank your ship"))
	_, err = sm.Delete([]int{id})
	if err == nil {
		t.Fatal("delete logic error")
	}
}

func TestSubjectManager_GetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := NewMockSubjectStore(mockCtrl)
	sm := SubjectManager{store: mockStore}

	sm.Init(list)

	tmp := sm.GetAll()
	if len(tmp) != len(list) {
		log.Println(tmp)
		log.Println(list)
		t.Fatal("get all failed")
	}
}

func TestSubjectManager_CheckSubjectList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := NewMockSubjectStore(mockCtrl)
	sm := SubjectManager{store: mockStore}

	sm.Init(list)

	load := []int{}
	for _, v := range list {
		load = append(load, v.ID)
	}
	if !sm.CheckSubjectList(load) {
		t.Fatal("check failed for existing id")
	}

	for k, v := range load {
		load[k] = v + 1
	}
	if sm.CheckSubjectList(load) {
		t.Fatal("check failed for non-existing id")
	}
}
