package manager

import (
	"github.com/arong/dean/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"testing"

	"github.com/arong/dean/models"
)

var students = models.StudentList{
	models.StudentInfo{Status: 1, StudentID: 1, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	models.StudentInfo{Status: 1, StudentID: 2, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	models.StudentInfo{Status: 1, StudentID: 3, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	models.StudentInfo{Status: 1, StudentID: 4, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	models.StudentInfo{Status: 1, StudentID: 5, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	models.StudentInfo{Status: 1, StudentID: 6, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	models.StudentInfo{Status: 1, StudentID: 7, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	models.StudentInfo{Status: 1, StudentID: 8, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	models.StudentInfo{Status: 1, StudentID: 9, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
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
				curr, err := s.GetUser(1)
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
			got, err := um.AddUser(tt.args.u)
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
