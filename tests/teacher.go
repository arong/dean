package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/arong/dean/base"
	"github.com/arong/dean/models"
)

func getTeacherList(host string) (models.TeacherList, error) {
	resp, err := sendPostRequest(host+"/api/v1/dean/teacher/filter", base.CommPage{Page: 1, Size: 20})
	if err != nil {
		return nil, err
	}

	ret := struct {
		Total int                `json:"total"`
		List  models.TeacherList `json:"list"`
	}{}

	err = json.Unmarshal(resp, &ret)
	for _, v := range ret.List {
		fmt.Println(v.Name)
	}
	return ret.List, nil
}

func getTeacherInfo(id int64) (*models.Teacher, error) {
	resp, err := sendGetRequest(host+"/api/v1/dean/teacher/info/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return nil, err
	}

	ret := &models.Teacher{}

	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func addTeacher() {

}

func modTeacher() {

}

func delTeacher(id int) {

}
