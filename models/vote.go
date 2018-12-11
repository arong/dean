package models

// VM a global handler
var VM voteManager

type voteManager struct {
	// teacherID -> votes
	votes map[int64]*ScoreInfo
}

func (vm *voteManager) Init() {
	vm.votes = make(map[int64]*ScoreInfo)
}

// AddScore calculate the score of a teacher
func (s *ScoreInfo) AddScore(score int) {
	s.votes = append(s.votes, score)
	total := 0.0
	for _, v := range s.votes {
		total += float64(v)
	}
	s.Average = total / float64(len(s.votes))
}

func (vm *voteManager) CastVote(votes []*VoteMeta) error {
	for _, v := range votes {
		if _, err := Tm.GetTeacherInfo(v.TeacherID); err == nil {
			val, ok := vm.votes[v.TeacherID]
			if !ok {
				val = &ScoreInfo{}
				vm.votes[v.TeacherID] = val
			}
			val.AddScore(v.Score)

		}
	}
	return nil
}

func (vm *voteManager) GetScore(teacherID int64) (*ScoreInfo, error) {
	ret := &ScoreInfo{}
	if val, ok := vm.votes[teacherID]; ok {
		*ret = *val
		return ret, nil
	}
	return ret, errNotExist
}

func (vm *voteManager) GetAll() []*ScoreInfo {
	ret := []*ScoreInfo{}
	for _, v := range vm.votes {
		tmp := *v
		ret = append(ret, &tmp)
	}
	return ret
}
