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
