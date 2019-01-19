package manager

import "github.com/arong/dean/models"

//SSM is global student score manager
var SSM StudentScoreManager

type StudentScoreManager struct {
	currentYear int // 当前学年
	currentTerm int // 当前学期
	currentExam int // 当前考试
	score       map[int64]models.YearScoreList
}

// AddRecord: add new record
func (ssm StudentScoreManager) AddRecord(r models.StudentScore) error {
	return nil
}

func (ssm StudentScoreManager) getStudentScore(studentID int64) (models.StudentScore, error) {
	item := models.StudentScore{}
	score, ok := ssm.score[studentID]
	if !ok {
		return item, errNotExist
	}
	for _, sy := range score {
		if sy.Year != ssm.currentYear {
			continue
		}
		for _, st := range sy.TermScores {
			if st.TermID != ssm.currentTerm {
				continue
			}
			item.TermID = st.TermID
			for _, se := range st.ExamsScores {
				if se.Exam != ssm.currentExam {
					continue
				}
				item.Exam = se.Exam
				for _, v := range se.Scores {
					tmp := v
					item.Scores = append(item.Scores, tmp)
				}
			}
		}
	}
	return item, nil
}

func (ssm StudentScoreManager) getClassScore(classID int) (models.StudentScoreList, error) {
	ret := models.StudentScoreList{}
	classInfo, err := Cm.GetInfo(classID)
	if err != nil {
		return ret, err
	}
	students := classInfo.StudentList

	ret = ssm.getCurrentScore(students)
	return ret, nil
}

func (ssm StudentScoreManager) getGradeScore(grade int) (models.StudentScoreList, error) {
	studentID, err := Um.getStudentList(grade)
	if err != nil {
		return nil, err
	}

	ret := ssm.getCurrentScore(studentID)

	return ret, nil
}

func (ssm StudentScoreManager) getCurrentScore(sid []int64) models.StudentScoreList {
	ret := models.StudentScoreList{}
	for _, s := range sid {
		item := models.StudentScore{StudentID: s}
		score, ok := ssm.score[s]
		if !ok {
			continue
		}
		for _, sy := range score {
			if sy.Year != ssm.currentYear {
				continue
			}
			for _, st := range sy.TermScores {
				if st.TermID != ssm.currentTerm {
					continue
				}
				item.TermID = st.TermID
				for _, se := range st.ExamsScores {
					if se.Exam != ssm.currentExam {
						continue
					}
					item.Exam = se.Exam
					for _, v := range se.Scores {
						tmp := v
						item.Scores = append(item.Scores, tmp)
					}
				}
			}
		}
		ret = append(ret, item)
	}

	return ret
}
