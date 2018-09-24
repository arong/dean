package models

// a global handler
var Vm voteManager

type voteManager struct {
	// teacherID -> votes
	votes map[int]*ScoreInfo
}

type VoteMeta struct {
	TeacherID int // 教师ID
	Score     int // 评分
}

type ScoreInfo struct {
	votes   []int
	Average float64 `json:"average"`
}

func (vm *voteManager) Init() {
	vm.votes = make(map[int]*ScoreInfo)
}

func (s *ScoreInfo) AddScore(score int) {
	s.votes = append(s.votes, score)

}
func (vm *voteManager) CastVote(votes []*VoteMeta) error {
	for _, v := range votes {
		if _, err := Tm.GetTeacherInfo(v.TeacherID); err == nil {
			if val, ok := vm.votes[v.TeacherID]; !ok {
				val = &ScoreInfo{}
				val.AddScore(v.Score)
				vm.votes[v.TeacherID] = val
			} else {
				val.AddScore(v.Score)
			}
		}
	}
	return nil
}

func (vm *voteManager) GetScore(teacherID int) (*ScoreInfo, error) {
	ret := &ScoreInfo{}
	if val, ok := vm.votes[teacherID]; ok {
		*ret = *val
		return ret, nil
	}
	return ret, ErrNotExist
}

func (vm *voteManager) GetAll() []*ScoreInfo {
	ret := []*ScoreInfo{}
	for _, v := range vm.votes {
		tmp := *v
		ret = append(ret, &tmp)
	}
	return ret
}
