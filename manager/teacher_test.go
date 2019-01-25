package manager

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/arong/dean/base"
	"github.com/arong/dean/mocks"
	"github.com/arong/dean/models"
	"github.com/golang/mock/gomock"
)

var teachers = models.TeacherList{
	{Status: base.StatusValid, TeacherMeta: models.TeacherMeta{TeacherID: 897, Name: "赵钱孙", Gender: eGenderMale, Mobile: "18621615580"}}, // standard
	{Status: base.StatusValid, TeacherMeta: models.TeacherMeta{TeacherID: 898, Name: "司空跃", Gender: eGenderMale, Mobile: "15235043583"}},
	{Status: base.StatusValid, TeacherMeta: models.TeacherMeta{TeacherID: 996, Name: "施圆梦", Gender: eGenderMale, Mobile: "18869111159"}},
	{Status: base.StatusValid, TeacherMeta: models.TeacherMeta{TeacherID: 998, Name: "蔡宜洋", Gender: eGenderMale, Mobile: "13355116391"}},
	{Status: base.StatusValid, TeacherMeta: models.TeacherMeta{TeacherID: 999, Name: "萧箔宸", Gender: eGenderMale, Mobile: "17843084173"}},
	{Status: base.StatusValid, TeacherMeta: models.TeacherMeta{TeacherID: 879, Name: "赵静姝", Gender: eGenderFemale, Mobile: "17788285999"}},
	{Status: base.StatusValid, TeacherMeta: models.TeacherMeta{TeacherID: 887, Name: "赵艳玲", Gender: eGenderFemale, Mobile: "18786322611"}},
	{Status: base.StatusValid, TeacherMeta: models.TeacherMeta{TeacherID: 857, Name: "赵炜旭", Gender: eGenderMale, Mobile: "13976381888"}},
	{Status: base.StatusValid, TeacherMeta: models.TeacherMeta{TeacherID: 497, Name: "赵锦明", Gender: eGenderMale, Mobile: "18169490121"}},
}

func TestTeacherManager_Init(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockTeacherStore(mockCtrl)
	tm := TeacherManager{store: mockStore}
	tm.Init(nil)

	if tm.list == nil {
		t.Fatal("bad behaviour")
	}
}

func TestTeacherManager_AddTeacher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockTeacherStore(mockCtrl)
	tm := TeacherManager{store: mockStore}
	tm.Init(nil)

	for k, v := range teachers {
		mockStore.EXPECT().SaveTeacher(gomock.Any()).Return(int64(k), nil)
		id, err := tm.AddTeacher(&v)
		if err != nil || id != int64(k) {
			t.Fatal("add teacher failed")
		}
	}

	for _, v := range teachers {
		_, err := tm.AddTeacher(&v)
		if err == nil {
			t.Error("add teacher failed")
		}
	}

	_, err := tm.AddTeacher(&models.Teacher{TeacherMeta: models.TeacherMeta{Name: strings.Repeat("烫", 64)}})
	if err == nil {
		t.Error("add teacher failed")
	}

	tm.Init(nil)
	var id int64
	mockStore.EXPECT().SaveTeacher(gomock.Any()).Return(id, errors.New("sank your ship"))
	tmp := teachers[0]
	_, err = tm.AddTeacher(&tmp)
	if err == nil {
		t.Error("add logic failed")
	}
}

func TestTeacherManager_ModTeacher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockTeacherStore(mockCtrl)

	tests := []struct {
		name    string
		list    models.TeacherList
		arg     models.Teacher
		wantErr bool
		before  func() // hooks to do some action before running test
	}{
		{
			name:    "invalid teacher id",
			list:    nil,
			arg:     models.Teacher{},
			wantErr: true,
		},
		{
			name:    "invalid input",
			list:    nil,
			arg:     models.Teacher{TeacherMeta: models.TeacherMeta{TeacherID: 1, Gender: 8}},
			wantErr: true,
		},
		{
			name:    "teacher id not exist",
			list:    nil,
			arg:     teachers[0],
			wantErr: true,
		},
		{
			name:    "update non-exist value",
			list:    teachers,
			arg:     teachers[0],
			wantErr: false,
		},
		{
			name:    "change name",
			list:    teachers,
			arg:     func() models.Teacher { tmp := teachers[0]; tmp.Name += "modify"; return tmp }(),
			wantErr: true,
		},
		{
			name: "normal change",
			list: teachers,
			arg: func() models.Teacher {
				tmp := teachers[1]
				tmp.Mobile = "11111111111"
				tmp.SubjectID = 89
				tmp.Address = "烫烫烫"
				tmp.Gender = eGenderFemale
				tmp.Birthday = "1987-09-12"
				return tmp
			}(),
			wantErr: false,
			before: func() {
				mockStore.EXPECT().UpdateTeacher(gomock.Any()).Return(nil)
			},
		},
		{
			name: "storage failed",
			list: teachers,
			arg: func() models.Teacher {
				tmp := teachers[0]
				tmp.Mobile = "11111111111"
				return tmp
			}(),
			wantErr: true,
			before: func() {
				mockStore.EXPECT().UpdateTeacher(gomock.Any()).Return(errors.New("mock storage failed"))
			},
		},
	}

	for _, v := range tests {
		tm := TeacherManager{store: mockStore}
		tm.Init(v.list)
		if v.before != nil {
			v.before()
		}
		err := tm.UpdateTeacher(&v.arg)
		if (err != nil) != v.wantErr {
			t.Errorf("execute '%s' failed for %v", v.name, v.arg)
		}
	}
}

func TestTeacherManager_DelTeacher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockTeacherStore(mockCtrl)
	tests := []struct {
		name    string
		list    models.TeacherList
		args    []int64
		want    []int64
		wantErr bool
		before  func()
	}{
		{
			name: "id not exist",
			list: nil,
			args: func() []int64 {
				load := []int64{}
				for _, v := range teachers {
					load = append(load, v.TeacherID)
				}
				return load
			}(),
			want: func() []int64 {
				load := []int64{}
				for _, v := range teachers {
					load = append(load, v.TeacherID)
				}
				return load
			}(),
			wantErr: false,
		},
		{
			name: "normal delete",
			list: teachers,
			args: func() []int64 {
				load := []int64{}
				for _, v := range teachers {
					load = append(load, v.TeacherID)
				}
				return load
			}(),
			want:    []int64{},
			wantErr: false,
			before: func() {
				mockStore.EXPECT().DeleteTeacher(gomock.Any()).Return(nil)
			},
		},
		{
			name: "storage failed",
			list: teachers,
			args: func() []int64 {
				load := []int64{}
				for _, v := range teachers {
					load = append(load, v.TeacherID)
				}
				return load
			}(),
			want: func() []int64 {
				load := []int64{}
				for _, v := range teachers {
					load = append(load, v.TeacherID)
				}
				return load
			}(),
			wantErr: false,
			before: func() {
				mockStore.EXPECT().DeleteTeacher(gomock.Any()).Return(errors.New("sank your ship"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TeacherManager{store: mockStore}
			tm.Init(tt.list)

			if tt.before != nil {
				tt.before()
			}
			got, err := tm.DelTeacher(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("TeacherManager.DelTeacher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TeacherManager.DelTeacher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTeacherManager_GetTeacherInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockTeacherStore(mockCtrl)

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		list    models.TeacherList
		args    args
		want    TeacherInfoResp
		wantErr bool
	}{
		{
			name:    "normal get",
			list:    teachers,
			args:    args{id: teachers[0].TeacherID},
			want:    TeacherInfoResp{TeacherMeta: teachers[0].TeacherMeta},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TeacherManager{store: mockStore}
			tm.Init(tt.list)

			got, err := tm.GetTeacherInfo(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TeacherManager.GetTeacherInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TeacherManager.GetTeacherInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTeacherManager_Filter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mocks.NewMockTeacherStore(mockCtrl)

	type args struct {
		f models.TeacherFilter
	}
	tests := []struct {
		name string
		list models.TeacherList
		args args
		want base.CommList
	}{
		{
			name: "normal test",
			list: teachers,
			args: args{f: models.TeacherFilter{CommPage: base.CommPage{Page: 1, Size: len(teachers)}}},
			want: base.CommList{Total: len(teachers), List: func() models.TeacherList {
				load := models.TeacherList{}
				for _, v := range teachers {
					tmp := v
					load = append(load, tmp)
				}
				sort.Sort(load)
				return load
			}()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TeacherManager{store: mockStore}
			tm.Init(tt.list)

			if got := tm.Filter(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TeacherManager.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTeacherManager_GetAll(t *testing.T) {
	tests := []struct {
		name string
		list models.TeacherList
		want base.CommList
	}{
		{
			name: "normal get",
			list: teachers,
			want: base.CommList{Total: len(teachers), List: func() simpleTeacherList {
				ret := simpleTeacherList{}
				for _, v := range teachers {
					ret = append(ret, simpleTeacher{Name: v.Name, ID: v.TeacherID})
				}
				return ret
			}()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TeacherManager{}
			tm.Init(tt.list)
			if got := tm.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TeacherManager.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTeacherManager_IsTeacherExist(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name string
		list models.TeacherList
		args args
		want bool
	}{
		{
			name: "normal test",
			list: teachers,
			args: args{id: teachers[0].TeacherID},
			want: true,
		},
		{
			name: "not exist",
			list: nil,
			args: args{id: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TeacherManager{}
			tm.Init(tt.list)
			if got := tm.IsTeacherExist(tt.args.id); got != tt.want {
				t.Errorf("TeacherManager.IsTeacherExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTeacherManager_CheckInstructorList(t *testing.T) {
	type args struct {
		list models.InstructorList
	}
	tests := []struct {
		name    string
		list    models.TeacherList
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			list: teachers,
			args: args{list: func() models.InstructorList {
				ret := models.InstructorList{}
				for _, v := range teachers {
					ret = append(ret, models.InstructorInfo{
						TeacherID: v.TeacherID,
						SubjectID: v.SubjectID,
					})
				}
				return ret
			}(),
			},
			wantErr: false,
		},
		{
			name: "subject not match",
			list: teachers,
			args: args{list: func() models.InstructorList {
				ret := models.InstructorList{}
				for _, v := range teachers {
					ret = append(ret, models.InstructorInfo{
						TeacherID: v.TeacherID,
						SubjectID: v.SubjectID + 1,
					})
				}
				return ret
			}(),
			},
			wantErr: true,
		},
		{
			name:    "teacher not exist",
			list:    nil,
			args:    args{list: models.InstructorList{models.InstructorInfo{TeacherID: 1, SubjectID: 1}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TeacherManager{}
			tm.Init(tt.list)
			if err := tm.CheckInstructorList(tt.args.list); (err != nil) != tt.wantErr {
				t.Errorf("TeacherManager.CheckInstructorList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTeacherManager_IsExist(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name string
		list models.TeacherList
		args args
		want bool
	}{
		{
			name: "normal test",
			list: teachers,
			args: args{id: teachers[0].TeacherID},
			want: true,
		},
		{
			name: "not exist",
			list: nil,
			args: args{id: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TeacherManager{}
			tm.Init(tt.list)
			if got := tm.IsExist(tt.args.id); got != tt.want {
				t.Errorf("TeacherManager.IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
