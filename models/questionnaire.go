package models

import "github.com/astaxie/beego/logs"

var Qm questionnaireManager

type questionnaireManager struct {
	questionnaires map[int]*QuestionnaireInfo
	titleMap       map[string]*QuestionnaireInfo
}

func (q *questionnaireManager) Add(info *QuestionnaireInfo) error {
	if _, ok := q.titleMap[info.Title]; ok {
		logs.Debug("[questionnaireManager::Add] name duplicated")
		return nil
	}

	err := Ma.AddQuestionnaire(info)
	if err != nil {
		logs.Error("[questionnaireManager::Add] failed")
		return err
	}

	q.questionnaires[info.QuestionnaireID] = info
	q.titleMap[info.Title] = info
	return nil
}

func (q *questionnaireManager) Update(info *QuestionnaireInfo) error {
	return nil
}

func (qm *questionnaireManager) Filter() (QuestionnaireList, error ){
	ret := QuestionnaireList{}
	for _, v := range qm.questionnaires {
		tmp := QuestionnaireInfo{
			QuestionnaireID: v.QuestionnaireID,
			Title:           v.Title,
			StartTime:       v.StartTime,
			StopTime:        v.StopTime,
			Status:          v.Status,
		}
		ret = append(ret, tmp)
	}
	return ret,nil
}
