package manager

import (
	"testing"

	"github.com/arong/dean/models"
)

func Test_userManager_Init(t *testing.T) {
	students := models.StudentList{
		models.StudentInfo{},
	}
	s := StudentManager{}
	s.Init(students)

	if len(s.list) != 0 {
		t.Error("init failed")
	}
}
