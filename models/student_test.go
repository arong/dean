package models

import (
	"strings"
	"testing"
)

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
