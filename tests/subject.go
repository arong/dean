package main

import (
	"encoding/json"
	"log"

	"github.com/arong/dean/models"
)

func GetAllSubjects() (models.SubjectList, error) {
	resp, err := sendGetRequest(host+"/api/v1/dean/subject/list", nil)
	if err != nil {
		return nil, err
	}

	ret := models.SubjectList{}

	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func AddMultiSubject() models.SubjectList {
	subject := models.SubjectList{
		{Key: "Chinese", Name: "语文"},
		{Key: "Math", Name: "数学"},
		{Key: "English", Name: "英语"},
		{Key: "Physics", Name: "物理"},
		{Key: "Chemistry", Name: "化学"},
		{Key: "Biology", Name: "生物"},
		{Key: "Politics", Name: "政治"},
		{Key: "History", Name: "历史"},
		{Key: "Geology", Name: "地理"},
		{Key: "P.E.", Name: "体育"},
	}

	for k, v := range subject {
		id, err := addSubject(v)
		if err != nil {
			log.Println("add subject failed")
		}
		subject[k].ID = id
	}
	return subject
}

func addSubject(info models.SubjectInfo) (int, error) {
	resp, err := sendPostRequest(host+"/api/v1/dean/subject/add", info)
	if err != nil {
		return 0, err
	}

	ret := struct {
		ID int `json:"id"`
	}{}

	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return 0, err
	}
	return ret.ID, nil
}

func DeleteAllSubject(list models.SubjectList) error {
	load := []int{}
	end := len(list) - 1

	for k, v := range list {
		load = append(load, v.ID)
		if (k+1)%100 == 0 || k == end {
			err := delSubject(load)
			if err != nil {
				return err
			}
			load = []int{}
		}
	}
	return nil
}

func delSubject(id []int) error {
	resp, err := sendPostRequest(host+"/api/v1/dean/subject/delete", struct {
		IDList []int `json:"id_list"`
	}{IDList: id})

	if err != nil {
		return err
	}

	if len(resp) > 0 {
		failedList := []int{}
		err = json.Unmarshal(resp, &failedList)
		if err != nil {
			return err
		}
		log.Println(failedList)
	}
	return nil
}
