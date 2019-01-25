package models

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/arong/dean/base"
)

var teachers = TeacherList{
	{Status: base.StatusValid, TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: eGenderMale, Mobile: "18621615580"}}, // standard
	{Status: base.StatusValid, TeacherMeta: TeacherMeta{Name: "司空跃", Gender: eGenderMale, Mobile: "15235043583"}},
	{Status: base.StatusValid, TeacherMeta: TeacherMeta{Name: "施圆梦", Gender: eGenderMale, Mobile: "18869111159"}},
	{Status: base.StatusValid, TeacherMeta: TeacherMeta{Name: "蔡宜洋", Gender: eGenderMale, Mobile: "13355116391"}},
	{Status: base.StatusValid, TeacherMeta: TeacherMeta{Name: "萧箔宸", Gender: eGenderMale, Mobile: "17843084173"}},
	{Status: base.StatusValid, TeacherMeta: TeacherMeta{Name: "赵静姝", Gender: eGenderFemale, Mobile: "17788285999"}},
	{Status: base.StatusValid, TeacherMeta: TeacherMeta{Name: "赵艳玲", Gender: eGenderFemale, Mobile: "18786322611"}},
	{Status: base.StatusValid, TeacherMeta: TeacherMeta{Name: "赵炜旭", Gender: eGenderMale, Mobile: "13976381888"}},
	{Status: base.StatusValid, TeacherMeta: TeacherMeta{Name: "赵锦明", Gender: eGenderMale, Mobile: "18169490121"}},
	{Status: base.StatusDeleted, TeacherMeta: TeacherMeta{Name: "赵明", Gender: eGenderMale, Mobile: "18169490121"}},
}

func TestTeacherMeta_Equal(t *testing.T) {
	prev := teachers[0]
	for k, v := range teachers {
		if !v.Equal(v.TeacherMeta) {
			t.Fatal("test equal failed")
		}

		if k != 0 && prev.Equal(v.TeacherMeta) {
			t.Error("should not equal")
		}
		prev = teachers[k]
	}
}

func TestTeacher_Check(t *testing.T) {
	in := []struct {
		t   Teacher
		err error
	}{
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: eGenderMale, Birthday: "1992-02-14", Address: "深圳市罗湖区"}}, err: nil}, // standard
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "", Gender: eGenderMale, Birthday: "1992-02-14", Address: "深圳市罗湖区"}}, err: ErrName},
		{t: Teacher{TeacherMeta: TeacherMeta{Name: strings.Repeat("烫", 17), Gender: eGenderMale, Birthday: "1992-02-14", Address: "深圳市罗湖区"}}, err: ErrName},
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: 0, Birthday: "1992-02-14", Address: "深圳市罗湖区"}}, err: ErrGender},
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: 4, Birthday: "1992-02-14", Address: "深圳市罗湖区"}}, err: ErrGender},
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: eGenderFemale, Birthday: "1992-02-14", Address: "深圳市罗湖区"}}, err: nil},
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: eGenderUnknown, Birthday: "1992-02-14", Address: "深圳市罗湖区"}}, err: nil},
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: eGenderMale, Address: strings.Repeat("烫", 65)}}, err: errAddress},
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: eGenderMale, Birthday: "1987-12-1 15:09:87"}}, err: ErrBirthday},
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: eGenderMale, Birthday: "0000-00-00"}}, err: ErrBirthday},
		{t: Teacher{TeacherMeta: TeacherMeta{Name: "赵钱孙", Gender: eGenderMale, Mobile: "127381983192"}}, err: errMobile},
	}

	for k, v := range in {
		err := v.t.Check()
		if err != v.err {
			t.Fatalf("%d check failed, err=%v", k, err)
		}
	}
}

func TestSort(t *testing.T) {
	list := TeacherList{}
	for _, v := range teachers {
		tmp := v
		tmp.TeacherID = int64(rand.Int())
		list = append(list, tmp)
	}
	sort.Sort(list)
	curr := list[0].TeacherID
	for _, v := range list[1:] {
		if curr > v.TeacherID {
			t.Error("sort failed")
		}
	}
}

func TestTeacherList_Page(t *testing.T) {
	total := len(teachers)

	tmp := teachers.Page(total, total)
	if len(tmp) > 0 {
		fmt.Println(tmp)
		t.Error("page failed")
	}

	tmp = teachers.Page(2, total-1)
	if len(tmp) != 1 {
		t.Error("page failed")
	}

	tmp = teachers.Page(1, total)
	if len(tmp) != total {
		t.Error("page failed")
	}
}

func TestTeacherList_Filter(t *testing.T) {
	type args struct {
		f TeacherFilter
	}
	tests := []struct {
		name  string
		tl    TeacherList
		args  args
		want  int
		want1 TeacherList
	}{
		{
			name:  "gender filter",
			tl:    teachers,
			args:  args{f: TeacherFilter{Gender: eGenderMale, CommPage: base.CommPage{1, len(teachers)}}},
			want:  7,
			want1: TeacherList{teachers[0], teachers[1], teachers[2], teachers[3], teachers[4], teachers[7], teachers[8]},
		},
		{
			name:  "filter name",
			tl:    teachers,
			args:  args{f: TeacherFilter{Name: "non-existing-name", CommPage: base.CommPage{Page: 1, Size: len(teachers)}}},
			want:  0,
			want1: TeacherList{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.tl.Filter(tt.args.f)
			if got != tt.want {
				t.Errorf("TeacherList.Filter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("TeacherList.Filter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
