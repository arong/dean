package models

import (
	"time"

	"github.com/arong/dean/base"
	"github.com/astaxie/beego/logs"
)

var QuestionnaireManager questionnaireManager

func init() {
	QuestionnaireManager.titleMap = make(map[string]*QuestionnaireInfo)
	QuestionnaireManager.questions = make(map[int]*QuestionInfo)
	QuestionnaireManager.score = make(map[int64]TeacherScore)
}

type questionnaireManager struct {
	questionnaires map[int]*QuestionnaireInfo
	titleMap       map[string]*QuestionnaireInfo
	questions      map[int]*QuestionInfo  // all question
	score          map[int64]TeacherScore // teacher score
	//page map[int]
}

func (q *questionnaireManager) Init(idMap map[int]*QuestionnaireInfo) {
	q.questionnaires = idMap
	for _, v := range idMap {
		q.titleMap[v.Title] = v
	}
}

func (q *questionnaireManager) Add(info *QuestionnaireInfo) (int, error) {
	if _, ok := q.titleMap[info.Title]; ok {
		logs.Debug("[questionnaireManager::Add] name duplicated")
		return 0, errExist
	}

	err := Ma.AddQuestionnaire(info)
	if err != nil {
		logs.Error("[questionnaireManager::Add] failed", "err", err)
		return 0, err
	}

	q.questionnaires[info.QuestionnaireID] = info
	q.titleMap[info.Title] = info
	return info.QuestionnaireID, nil
}

// many thing to do in update
func (q *questionnaireManager) Update(info *QuestionnaireInfo) error {
	curr, ok := q.questionnaires[info.QuestionnaireID]
	if !ok {
		return errNotExist
	}

	if curr.Status != QStatusDraft {
		return errPermission
	}

	if curr.Equal(*info) {
		logs.Info("[questionnaireManager::Update] nothing to do")
		return nil
	}

	backup := *curr
	if curr.Title != info.Title {
		curr.Title = info.Title
	}

	if curr.StartTime != info.StartTime {
		tmp, err := time.Parse(base.DateTimeFormat, info.StartTime)
		if err != nil {

		}
		curr.startTime = tmp
		curr.StartTime = info.StartTime
	}

	if curr.StopTime != info.StopTime {
		tmp, err := time.Parse(base.DateTimeFormat, info.StopTime)
		if err != nil {
			return err
		}
		curr.stopTime = tmp
		curr.StopTime = info.StopTime
	}

	err := Ma.UpdateQuestionnaire(curr)
	if err != nil {
		curr = &backup
		logs.Info("[questionnaireManager::Update] UpdateQuestionnaire failed", "err", err)
		return err
	}
	return nil
}

type GenRequest struct {
	StudentID       int64 `json:"-"`
	QuestionnaireID int   `json:"questionnaire_id"`
}

func (qm *questionnaireManager) Generate(request GenRequest) (SurveyPages, error) {
	q, ok := qm.questionnaires[request.QuestionnaireID]
	if !ok {
		return nil, errNotExist
	}

	if q.Status != QStatusPublished {
		logs.Debug("[questionnaireManager::Generate] status not allowed")
		return nil, errNotExist
	}

	if len(q.Questions) == 0 {
		logs.Debug("[questionnaireManager::Generate] no question")
		return nil, errNotExist
	}

	student, err := Um.GetUser(request.StudentID)
	if err != nil {
		logs.Error("[questionnaireManager::Generate] GetUser failed", "err", err)
		return nil, err
	}

	class, err := Cm.GetInfo(student.ClassID)
	if err != nil {
		logs.Error("[questionnaireManager::Generate] GetInfo failed", "err", err)
		return nil, err
	}

	logs.Debug("[questionnaireManager::Generate] q", class.TeacherList)
	survey := SurveyPages{}
	for _, v := range class.TeacherList {
		page := SurveyPage{TeacherID: v.TeacherID, TeacherName: v.Teacher}
		for _, val := range q.Questions {
			if len(val.Scope) > 0 {
				found := false
				for _, i := range val.Scope {
					if i == v.SubjectID {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			page.Questions = append(page.Questions, val)
		}
		survey = append(survey, page)
	}

	return survey, nil
}

func (q *questionnaireManager) Delete(id int) error {
	curr, ok := q.questionnaires[id]
	if !ok {
		return errNotExist
	}

	if curr.Status != QStatusDraft {
		return errPermission
	}

	err := Ma.DeleteQuestionnaire(id)
	if err != nil {
		return err
	}

	delete(q.titleMap, curr.Title)
	delete(q.questionnaires, id)
	return nil
}

func (qm *questionnaireManager) Filter() (QuestionnaireList, error) {
	ret := QuestionnaireList{}
	for _, v := range qm.questionnaires {
		tmp := QuestionnaireInfo{
			QuestionnaireID: v.QuestionnaireID,
			Title:           v.Title,
			StartTime:       v.StartTime,
			StopTime:        v.StopTime,
			Status:          v.Status,
			Editor:          v.Editor,
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

// AddQuestion add question to questionnaire
func (qm *questionnaireManager) AddQuestion(info *QuestionInfo) (int, error) {
	q, ok := qm.questionnaires[info.QuestionnaireID]
	if !ok {
		return 0, errNotExist
	}

	if q.Status != QStatusDraft {
		return 0, errPermission
	}

	for _, v := range q.Questions {
		if v.Question == info.Question {
			return 0, errExist
		}
	}

	// insert to database
	var err error
	info.QuestionID, err = Ma.AddQuestion(info.QuestionnaireID, info)
	if err != nil {
		logs.Warn("[AddQuestion] AddQuestion failed", "err", err)
		return 0, err
	}

	qm.questions[info.QuestionID] = info
	q.Questions = append(q.Questions, info)
	return info.QuestionID, nil
}

// UpdateQuestion add question to questionnaire
func (qm *questionnaireManager) UpdateQuestion(info *QuestionInfo) error {
	curr, ok := qm.questions[info.QuestionID]
	if !ok {
		return errNotExist
	}

	q, ok := qm.questionnaires[curr.QuestionnaireID]
	if !ok {
		return errNotExist
	}

	if q.Status != QStatusDraft {
		return errPermission
	}

	if curr.Equal(info) {
		logs.Debug("[questionnaireManager::UpdateQuestion] nothing to do")
		return nil
	}

	for _, v := range q.Questions {
		if v.Question == info.Question && v.QuestionID != info.QuestionID {
			return errExist
		}
	}

	backup := *curr
	curr.Index = info.Index
	curr.Type = info.Type
	curr.Required = info.Required
	curr.Question = info.Question
	curr.Options = info.Options.FilterEmpty()

	// insert to database
	var err error
	info.QuestionID, err = Ma.UpdateQuestion(curr)
	if err != nil {
		logs.Warn("[AddQuestion] AddQuestion failed", "err", err)
		curr = &backup
		return err
	}

	return nil
}

func (qm *questionnaireManager) DeleteQuestion(id int) error {
	curr, ok := qm.questions[id]
	if !ok {
		return errNotExist
	}

	q, ok := qm.questionnaires[curr.QuestionnaireID]
	if !ok {
		return errNotExist
	}

	if q.Status != QStatusDraft {
		return errPermission
	}

	err := Ma.DeleteQuestion(id)
	if err != nil {
		logs.Info("[questionnaireManager::DeleteQuestion] DeleteQuestion failed", "err", err)
		return err
	}

	delete(qm.questions, id)
	if len(q.Questions) > 0 {
		list := QuestionList{}
		for _, v := range q.Questions {
			if v.QuestionID == id {
				continue
			}
			list = append(list, v)
		}
		q.Questions = list
	}

	return nil
}

func (qm *questionnaireManager) GetQuestionInfo(id int) (*QuestionInfo, error) {
	curr, ok := qm.questions[id]
	if !ok {
		return nil, errNotExist
	}
	return curr, nil
}

func (qm *questionnaireManager) GetQuestions(questionnaireID int) (QuestionList, error) {
	curr, ok := qm.questionnaires[questionnaireID]
	if !ok {
		return nil, errNotExist
	}

	list := QuestionList{}
	for _, v := range curr.Questions {
		list = append(list, v)
	}
	return list, nil
}

// submit questionnaire
//Submit submit questionnaire of a student
func (qm *questionnaireManager) Submit(req QuestionnaireSubmit) error {
	// get question info
	curr, ok := qm.questionnaires[req.QuestionnaireID]
	if !ok {
		logs.Debug("[questionnaireManager::Submit] questionnaire not found")
		return errNotExist
	}

	if curr.Status != QStatusPublished {
		return errPermission
	}

	studentInfo, err := Um.GetUser(req.StudentID)
	if err != nil {
		return err
	}

	classInfo, err := Cm.GetInfo(studentInfo.ClassID)
	if err != nil {
		return err
	}

	// check to see if student have right to vote to the teacher
	teacher := make(map[int64]int)
	for _, v := range classInfo.TeacherList {
		teacher[v.TeacherID] = v.SubjectID
	}

	for _, v := range req.TeacherAnswers {
		if _, ok := teacher[v.TeacherID]; !ok {
			return errPermission
		}

		ans := make(map[int]*AnswerInfo)
		for _, val := range v.Answers {
			ans[val.QuestionID] = val
		}

		for _, question := range curr.Questions {
			tmp, ok := ans[question.QuestionID]
			// check required question
			if question.Required && !ok {
				return errNotExist
			}

			if len(question.Scope) > 0 {
				found := false
				sid := teacher[v.TeacherID]
				for _, v := range question.Scope {
					if v == sid {
						found = true
						break
					}
				}
				if !found {
					return errPermission
				}
			}

			switch question.Type {
			case QuestionTypeSingleSelection:
				tScore := qm.score[v.TeacherID]
				if w, ok := tmp.Answer.(float64); ok {
					choice := int(w)
					tScore.Count++
					tScore.Meta[choice] = append(tScore.Meta[choice], sourceMeta{Grade: classInfo.Grade, Index: classInfo.Index})
				} else {
					return errInvalidInput
				}
			case QuestionTypeMultiSelection:
				tScore := qm.score[v.TeacherID]
				if w, ok := tmp.Answer.([]float64); ok {
					for _, i := range w {
						choice := int(i)
						tScore.Count++
						tScore.Meta[choice] = append(tScore.Meta[choice], sourceMeta{Grade: classInfo.Grade, Index: classInfo.Index})
					}
				} else {
					return errInvalidInput
				}
			case QuestionTypeText:
				tScore := qm.score[v.TeacherID]
				if w, ok := tmp.Answer.(string); ok {
					if w != "" {
						tScore.Remark = append(tScore.Remark, w)
					}
				} else {
					return errInvalidInput
				}
			default:
				logs.Error("[questionnaireManager::Submit] internal bug found")
			}
		}
	}

	return nil
}
