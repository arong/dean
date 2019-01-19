package models

type StudentInfo struct {
	Age        int    `json:"age"`
	Gender     int    `json:"gender"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	Address    string `json:"address"`
	Birthday   string `json:"birthday"`
	StudentID  int64  `json:"student_id"`
	Status     int    `json:"status"`
	ClassID    int    `json:"class_id"`
	RegisterID string `json:"register_id"` // 学号
}

type StudentList []StudentInfo

func (s StudentList) Filter(f StudentFilter) StudentList {
	ret := StudentList{}
	list := IntList{}

	for k, v := range s {
		if f.Name != "" && f.Name != v.Name {
			continue
		}

		if f.Number != "" && f.Number != v.RegisterID {
			continue
		}

		list = append(list, k)
	}
	list = list.Page(f.CommPage)
	for _, v := range list {
		ret = append(ret, s[v])
	}
	return ret
}

func (cl StudentList) Len() int {
	return len(cl)
}

func (cl StudentList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

func (cl StudentList) Less(i, j int) bool {
	return cl[i].StudentID < cl[j].StudentID
}

type StudentStore interface {
	SaveStudent(StudentInfo) (int64, error)
	UpdateStudent(info StudentInfo) error
	DeleteStudent([]int64) error
}
