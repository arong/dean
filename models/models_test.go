package models

import (
	"testing"
)

func TestQuestionnaireInfo_IsSame(t *testing.T) {
	l := QuestionnaireInfo{
		QuestionnaireID: 1,
		Title:           "l",
		Questions:       QuestionList{&QuestionInfo{Question: "how are you"}},
	}

	r := QuestionnaireInfo{
		QuestionnaireID: 1,
		Title:           "r",
		Questions:       QuestionList{},
	}

	if l.Equal(r) {
		t.Error("logic break")
	}
}

//func TestTeacher_Check(t *testing.T) {
//	in := []struct {
//		t   Teacher
//		err error
//	}{
//		{t: Teacher{Name: "赵钱孙", Gender: eGenderMale, Birthday: "1992-02-14", Address: "深圳市罗湖区"}, err: nil}, // standard
//		{t: Teacher{Name: "", Gender: eGenderMale, Birthday: "1992-02-14", Address: "深圳市罗湖区"}, err: ErrName},
//		{t: Teacher{Name: strings.Repeat("烫", 17), Gender: eGenderMale, Birthday: "1992-02-14", Address: "深圳市罗湖区"}, err: ErrName},
//		{t: Teacher{Name: "赵钱孙", Gender: 0, Birthday: "1992-02-14", Address: "深圳市罗湖区"}, err: ErrGender},
//		{t: Teacher{Name: "赵钱孙", Gender: 4, Birthday: "1992-02-14", Address: "深圳市罗湖区"}, err: ErrGender},
//		{t: Teacher{Name: "赵钱孙", Gender: eGenderFemale, Birthday: "1992-02-14", Address: "深圳市罗湖区"}, err: nil},
//		{t: Teacher{Name: "赵钱孙", Gender: eGenderUnknown, Birthday: "1992-02-14", Address: "深圳市罗湖区"}, err: nil},
//		{t: Teacher{Name: "赵钱孙", Gender: eGenderMale, Address: strings.Repeat("烫", 65)}, err: errAddress},
//		{t: Teacher{Name: "赵钱孙", Gender: eGenderMale, Birthday: "1987-12-1 15:09:87"}, err: ErrBirthday},
//		{t: Teacher{Name: "赵钱孙", Gender: eGenderMale}, err: nil},
//	}
//
//	for k, v := range in {
//		err := v.t.Check()
//		if err != v.err {
//			t.Fatalf("%d check failed, err=%v", k, err)
//		}
//	}
//}
