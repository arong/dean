package models

import (
	"reflect"
	"testing"

	"github.com/arong/dean/base"
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

func TestIntList_Page(t *testing.T) {
	type args struct {
		page base.CommPage
	}
	tests := []struct {
		name string
		il   IntList
		args args
		want IntList
	}{
		{
			name: "normal range",
			il: func() IntList {
				ret := IntList{}
				for i := 0; i < 1000; i++ {
					ret = append(ret, i)
				}
				return ret
			}(),
			args: args{page: base.CommPage{Page: 1, Size: 10}},
			want: func() IntList {
				ret := IntList{}
				for i := 0; i < 10; i++ {
					ret = append(ret, i)
				}
				return ret
			}(),
		},
		{
			name: "out of range",
			il:   IntList{},
			args: args{page: base.CommPage{Page: 1, Size: 10}},
			want: IntList{},
		},
		{
			name: "out range",
			il: func() IntList {
				ret := IntList{}
				for i := 0; i < 1000; i++ {
					ret = append(ret, i)
				}
				return ret
			}(),
			args: args{page: base.CommPage{Page: 2, Size: 999}},
			want: IntList{999},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.il.Page(tt.args.page); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntList.Page() = %v, want %v", got, tt.want)
			}
		})
	}
}
