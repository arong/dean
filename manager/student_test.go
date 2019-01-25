package manager

import (
	"reflect"
	"testing"

	"github.com/arong/dean/base"
	"github.com/arong/dean/mocks"
	"github.com/arong/dean/models"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

var students = models.StudentList{
	models.StudentInfo{Status: 1, StudentID: 1, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190101"},
	models.StudentInfo{Status: 1, StudentID: 2, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	models.StudentInfo{Status: 1, StudentID: 3, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190103"},
	models.StudentInfo{Status: 1, StudentID: 4, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190104"},
	models.StudentInfo{Status: 1, StudentID: 5, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190105"},
	models.StudentInfo{Status: 1, StudentID: 6, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190106"},
	models.StudentInfo{Status: 1, StudentID: 7, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190107"},
	models.StudentInfo{Status: 1, StudentID: 8, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190108"},
	models.StudentInfo{Status: 1, StudentID: 9, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190109"},
}

func Test_userManager_Init(t *testing.T) {
	students := models.StudentList{
		models.StudentInfo{Gender: eGenderFemale, Name: "Jasmine", ClassID: 1},
	}
	s := StudentManager{}
	s.Init(students)

	if len(s.list) != 0 {
		t.Error("init failed")
	}
}

func TestStudentManager_AddUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockStudentStore(mockCtrl)

	type args struct {
		u models.StudentInfo
	}
	tests := []struct {
		name    string
		list    models.StudentList
		args    args
		want    int64
		wantErr bool
		before  func()
		after   func(*StudentManager)
	}{
		{
			name:    "add duplicated value",
			list:    students,
			args:    args{u: students[0]},
			want:    0,
			wantErr: true,
		},
		{
			name:    "normal add",
			list:    nil,
			args:    args{u: students[0]},
			want:    1,
			wantErr: false,
			before: func() {
				mockStore.EXPECT().SaveStudent(gomock.Any()).Return(int64(1), nil)
			},
			after: func(s *StudentManager) {
				curr, err := s.GetStudent(1)
				if err != nil || !curr.Equal(students[0]) {
					t.Error("add failed")
				}
			},
		},
		{
			name:    "storage failed",
			list:    nil,
			args:    args{u: students[0]},
			want:    0,
			wantErr: true,
			before: func() {
				mockStore.EXPECT().SaveStudent(gomock.Any()).Return(int64(0), errors.New("sank your ship"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := &StudentManager{store: mockStore}
			um.Init(tt.list)
			if tt.before != nil {
				tt.before()
			}
			got, err := um.AddStudent(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudentManager.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("StudentManager.AddUser() = %v, want %v", got, tt.want)
			}

			if tt.after != nil {
				tt.after(um)
			}
		})
	}
}

func TestStudentManager_UpdateStudent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockStudentStore(mockCtrl)
	type args struct {
		u models.StudentInfo
	}
	errStore := errors.New("sank your ship")
	tests := []struct {
		name   string
		list   models.StudentList
		args   args
		err    error
		before func()
	}{
		{
			name: "normal test",
			list: models.StudentList{{Status: 1, StudentID: 1, Mobile: "18719876543"}},
			args: args{u: models.StudentInfo{StudentID: 1, Gender: eGenderFemale}},
			err:  nil,
			before: func() {
				mockStore.EXPECT().UpdateStudent(gomock.Any()).Return(nil)
			},
		},
		{
			name: "store failed",
			list: models.StudentList{{Status: 1, StudentID: 1, Mobile: "18719876543"}},
			args: args{u: models.StudentInfo{StudentID: 1, Gender: eGenderFemale}},
			err:  errStore,
			before: func() {
				mockStore.EXPECT().UpdateStudent(gomock.Any()).Return(errStore)
			},
		},
		{
			name: "invalid student id",
			list: nil,
			args: args{u: models.StudentInfo{StudentID: 1}},
			err:  errNotExist,
		},
		{
			name: "no change",
			list: students,
			args: args{u: students[0]},
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := &StudentManager{store: mockStore}
			um.Init(tt.list)

			if tt.before != nil {
				tt.before()
			}

			if err := um.UpdateStudent(tt.args.u); err != tt.err {
				t.Errorf("StudentManager.UpdateStudent() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestStudentManager_DelUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockStudentStore(mockCtrl)

	type args struct {
		uidList []int64
	}
	tests := []struct {
		name    string
		list    models.StudentList
		args    args
		want    []int64
		wantErr bool
		before  func()
		after   func(manager StudentManager)
	}{
		{
			name: "normal delete",
			list: students,
			args: args{uidList: func() []int64 {
				ret := make([]int64, 0, len(students))
				for _, v := range students {
					ret = append(ret, v.StudentID)
				}
				return ret
			}()},
			before: func() {
				mockStore.EXPECT().DeleteStudent(gomock.Any()).Return(nil)
			},
			want:    []int64{},
			wantErr: false,
			after: func(manager StudentManager) {
				if len(manager.idMap) != 0 || len(manager.uuidMap) != 0 {
					t.Error("delete failed")
				}
			},
		},
		{
			name:    "delete non-existing id",
			list:    nil,
			args:    args{uidList: []int64{1, 2, 3, 4, 5}},
			want:    []int64{1, 2, 3, 4, 5},
			wantErr: true,
		},
		{
			name:    "storage failed",
			list:    students,
			args:    args{uidList: []int64{1}},
			want:    []int64{1},
			wantErr: true,
			before: func() {
				mockStore.EXPECT().DeleteStudent(gomock.Any()).Return(errors.New("sank your ship"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := &StudentManager{store: mockStore}
			um.Init(tt.list)
			if tt.before != nil {
				tt.before()
			}
			got, err := um.DelStudent(tt.args.uidList)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudentManager.DelUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StudentManager.DelUser() = %v, want %v", got, tt.want)
			}
			if tt.after != nil {
				tt.after(*um)
			}
		})
	}
}

func TestStudentManager_Filter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockStudentStore(mockCtrl)

	type args struct {
		f models.StudentFilter
	}

	tests := []struct {
		name   string
		list   models.StudentList
		args   args
		want   base.CommList
		before func(*StudentManager)
	}{
		{
			name: "normal list",
			list: students,
			args: args{f: models.StudentFilter{CommPage: base.CommPage{Page: 1, Size: len(students)}}},
			want: base.CommList{Total: len(students), List: students},
		},
		{
			name: "after delete all",
			list: students,
			args: args{f: models.StudentFilter{CommPage: base.CommPage{Page: 1, Size: len(students)}}},
			want: base.CommList{Total: 0, List: models.StudentList{}},
			before: func(manager *StudentManager) {
				ids := []int64{}
				for _, v := range students {
					ids = append(ids, v.StudentID)
				}
				mockStore.EXPECT().DeleteStudent(gomock.Any()).Return(nil)
				manager.DelStudent(ids)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := &StudentManager{store: mockStore}
			um.Init(tt.list)

			if tt.before != nil {
				tt.before(um)
			}
			if got := um.Filter(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StudentManager.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStudentManager_GetStudentByRegisterNumber(t *testing.T) {
	type args struct {
		reg string
	}
	tests := []struct {
		name    string
		list    models.StudentList
		args    args
		want    models.StudentInfo
		wantErr bool
	}{
		{
			name:    "not exist",
			list:    students,
			args:    args{reg: "not exist"},
			want:    models.StudentInfo{},
			wantErr: true,
		},
		{
			name:    "normal",
			list:    students,
			args:    args{reg: students[0].RegisterID},
			want:    students[0],
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := &StudentManager{}
			um.Init(tt.list)
			got, err := um.GetStudentByRegisterNumber(tt.args.reg)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudentManager.GetStudentByRegisterNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StudentManager.GetStudentByRegisterNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStudentManager_GetStudent(t *testing.T) {
	type args struct {
		uid int64
	}
	tests := []struct {
		name    string
		list    models.StudentList
		args    args
		want    models.StudentInfo
		wantErr bool
	}{
		{
			name:    "normal get",
			list:    students,
			args:    args{uid: 1},
			want:    students[0],
			wantErr: false,
		},
		{
			name:    "get non-exist",
			list:    nil,
			args:    args{uid: 1},
			want:    models.StudentInfo{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := &StudentManager{}
			um.Init(tt.list)
			got, err := um.GetStudent(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudentManager.GetStudent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StudentManager.GetStudent() = %v, want %v", got, tt.want)
			}
		})
	}
}
