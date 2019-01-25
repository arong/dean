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

func TestIntList_RemoveZeroNegative(t *testing.T) {
	tests := []struct {
		name string
		il   IntList
		want IntList
	}{
		{
			name: "normal test",
			il:   IntList{1, 2, 4, 0, 0, 0, 0, 0, 5},
			want: IntList{1, 2, 4, 5},
		},
		{
			name: "all zero",
			il:   IntList{0, 0, 0, 0, 0, 0, 0, 0, 0},
			want: IntList{},
		},
		{
			name: "end zero",
			il:   IntList{1, 0, 0, 0, 0, 0, 0, 0, 0},
			want: IntList{1},
		},
		{
			name: "begin zero",
			il:   IntList{0, 0, 0, 0, 0, 0, 0, 0, 5},
			want: IntList{5},
		},
		{
			name: "no zero",
			il:   IntList{1, 2, 3, 4, 5, 6, 7, 8, 9},
			want: IntList{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.il.RemoveZeroNegative(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntList.RemoveZeroNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64List_RemoveZeroNegative(t *testing.T) {
	tests := []struct {
		name string
		il   Int64List
		want Int64List
	}{
		{
			name: "normal test",
			il:   Int64List{1, 2, 4, 0, 0, 0, 0, 0, 5},
			want: Int64List{1, 2, 4, 5},
		},
		{
			name: "all zero",
			il:   Int64List{0, 0, 0, 0, 0, 0, 0, 0, 0},
			want: Int64List{},
		},
		{
			name: "end zero",
			il:   Int64List{1, 0, 0, 0, 0, 0, 0, 0, 0},
			want: Int64List{1},
		},
		{
			name: "begin zero",
			il:   Int64List{0, 0, 0, 0, 0, 0, 0, 0, 5},
			want: Int64List{5},
		},
		{
			name: "no zero",
			il:   Int64List{1, 2, 3, 4, 5, 6, 7, 8, 9},
			want: Int64List{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.il.RemoveZeroNegative(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64List.RemoveZeroNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}
