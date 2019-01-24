package models

import (
	"reflect"
	"strings"
	"testing"

	"github.com/arong/dean/base"
)

var students = StudentList{
	StudentInfo{Status: 1, StudentID: 1, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190101"},
	StudentInfo{Status: 1, StudentID: 2, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190102"},
	StudentInfo{Status: 1, StudentID: 3, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190103"},
	StudentInfo{Status: 1, StudentID: 4, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190104"},
	StudentInfo{Status: 1, StudentID: 5, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190105"},
	StudentInfo{Status: 1, StudentID: 6, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190106"},
	StudentInfo{Status: 1, StudentID: 7, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190107"},
	StudentInfo{Status: 1, StudentID: 8, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190108"},
	StudentInfo{Status: 1, StudentID: 9, Name: "Jasmine", Gender: eGenderFemale, RegisterID: "20190109"},
}

func TestStudentInfo_Check(t *testing.T) {
	tests := []struct {
		name string
		args StudentInfo
		err  error
	}{
		{
			name: "check gender",
			args: StudentInfo{Name: "Jasmine", Mobile: "18765431234", Address: "Shanghai"},
			err:  ErrGender,
		},
		{
			name: "check name",
			args: StudentInfo{Gender: eGenderFemale, Name: "", Mobile: "18765431234", Address: "Shanghai"},
			err:  ErrName,
		},
		{
			name: "check mobile",
			args: StudentInfo{Gender: eGenderFemale, Name: "Jasmine", Mobile: "123456789001", Address: "Shanghai"},
			err:  errMobile,
		},
		{
			name: "check address",
			args: StudentInfo{Gender: eGenderFemale, Name: "Jasmine", Mobile: "18765431234", Address: strings.Repeat("Shanghai", 10)},
			err:  errAddress,
		},
		{
			name: "check birthday",
			args: StudentInfo{Gender: eGenderFemale, Name: "Jasmine", Mobile: "18765431234", Address: "Shanghai", Birthday: "1289131893712893712893"},
			err:  ErrBirthday,
		},
		{
			name: "check birthday",
			args: StudentInfo{Gender: eGenderFemale, Name: "Jasmine", Mobile: "18765431234", Address: "Shanghai", Birthday: "1987-98-01"},
			err:  ErrBirthday,
		},
		{
			name: "check birthday",
			args: StudentInfo{Gender: eGenderFemale, Name: "Jasmine", Mobile: "18765431234", Address: "Shanghai", Birthday: "1987-09-01"},
			err:  errRegisterID,
		},
		{
			name: "check birthday",
			args: StudentInfo{Gender: eGenderFemale, Name: "Jasmine", Mobile: "18765431234", Address: "Shanghai", Birthday: "1987-09-01", RegisterID: "foo"},
			err:  errRegisterID,
		},
		{
			name: "check birthday",
			args: StudentInfo{Gender: eGenderFemale, Name: "Jasmine", Mobile: "18765431234", Address: "Shanghai", Birthday: "1987-09-01", RegisterID: "20190109"},
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.Check(); err != tt.err {
				t.Errorf("StudentInfo.Check() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestStudentList_Filter(t *testing.T) {
	type args struct {
		f StudentFilter
	}
	tests := []struct {
		name  string
		s     StudentList
		args  args
		want  int
		want1 StudentList
	}{
		{
			name:  "normal test",
			s:     students,
			args:  args{f: StudentFilter{CommPage: base.CommPage{Page: 1, Size: len(students)}}},
			want:  len(students),
			want1: students,
		},
		{
			name:  "out -> out",
			s:     students,
			args:  args{f: StudentFilter{CommPage: base.CommPage{Page: 2, Size: len(students)}}},
			want:  len(students),
			want1: StudentList{},
		},
		{
			name:  "in -> out",
			s:     students,
			args:  args{f: StudentFilter{CommPage: base.CommPage{Page: 2, Size: len(students) - 1}}},
			want:  len(students),
			want1: StudentList{students[len(students)-1]},
		},
		{
			name: "in -> in",
			s:    students,
			args: args{f: StudentFilter{CommPage: base.CommPage{Page: 2, Size: 4}}},
			want: len(students),
			want1: func() StudentList {
				return students[4:8]
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.Filter(tt.args.f)
			if got != tt.want {
				t.Errorf("StudentList.Filter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("StudentList.Filter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
