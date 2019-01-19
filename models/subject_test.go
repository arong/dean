package models

import (
	"sort"
	"testing"

	"github.com/arong/dean/base"
)

var subject = SubjectList{
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

func TestSubjectInfo_Equal(t *testing.T) {
	for _, v := range subject {
		if !v.Equal(v) {
			t.Error("logic error")
		}
	}
}

func TestSubjectList_Sort(t *testing.T) {
	list := SubjectList{}
	for i := len(subject) - 1; i > -1; i-- {
		tmp := subject[i]
		list = append(list, tmp)
	}
	sort.Sort(list)
	curr := list[0]
	for _, v := range list[1:] {
		if curr.ID > v.ID {
			t.Error("sort failed")
		}
	}
}

func TestSubjectInfo_Check(t *testing.T) {
	type fields struct {
		Status int
		ID     int
		Name   string
		Key    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				Name: "语文",
				Key:  "chinese",
			},
			wantErr: false,
		},
		{
			name: "abnormal",
			fields: fields{
				Name: "语文",
				Key:  "语文",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SubjectInfo{
				Status: tt.fields.Status,
				ID:     tt.fields.ID,
				Name:   tt.fields.Name,
				Key:    tt.fields.Key,
			}
			if err := s.Check(); (err != nil) != tt.wantErr {
				t.Errorf("SubjectInfo.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
