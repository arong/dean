package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"

	"github.com/pkg/errors"

	"github.com/go-vgo/gt/http"

	"github.com/arong/dean/base"
	"github.com/arong/dean/models"
)

func GetTeacherListAll(host string) (models.TeacherList, error) {
	i := 1
	size := 20
	all := models.TeacherList{}
	uniq := make(map[int64]bool)
	for {
		curr, err := getTeacherList(host, i, size)
		if err != nil {
			return all, err
		}

		for _, v := range curr {
			_, ok := uniq[v.TeacherID]
			if ok {
				log.Println("bug found", "i", i)
			} else {
				uniq[v.TeacherID] = true
			}
		}
		all = append(all, curr...)
		if len(curr) < size {
			break
		}
		i++
	}
	log.Println("total uniq teacher", len(uniq))
	if _, ok := uniq[0]; ok {
		log.Println("garbage found")
	}
	return all, nil
}

func getTeacherList(host string, page, size int) (models.TeacherList, error) {
	resp, err := sendPostRequest(host+"/api/v1/dean/teacher/filter", base.CommPage{Page: page, Size: size})
	if err != nil {
		log.Println("sendPostRequest failed", err)
		return nil, err
	}

	ret := struct {
		Total int                `json:"total"`
		List  models.TeacherList `json:"list"`
	}{}

	err = json.Unmarshal(resp, &ret)
	if err != nil {
		log.Println("unrecognized return value", err)
		return nil, err
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

func AddMultiTeachers() error {
	// get basic info
	for i := 0; i < 10000; i++ {
		t := getFakeInfo()
		t.Gender = rand.Int()%2 + 1
		id, err := addTeacher(t)
		if err != nil {
			log.Println("add teacher failed", t)
			continue
		}

		if id == 0 {
			log.Println("invalid teacher id")
		}
	}
	return nil
}

func getFakeInfo() models.Teacher {
	resp, err := http.Get("http://192.168.231.132:8080/api", nil)
	if err != nil {
		log.Println("random source failed, check it")
		return models.Teacher{}
	}

	ret := struct {
		Name   string
		Mobile string
		IdNo   string
		Bank   string
		Email  string
		Addr   string
	}{}
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		log.Println("unmarshal failed")
		return models.Teacher{}
	}
	return models.Teacher{
		TeacherMeta: models.TeacherMeta{
			Name:    ret.Name,
			Mobile:  ret.Mobile,
			Address: ret.Addr,
		},
	}
}

func addTeacher(teacher models.Teacher) (int64, error) {
	resp, err := sendPostRequest(host+"/api/v1/dean/teacher/add", teacher)
	if err != nil {
		log.Println("add teacher failed", err)
		return 0, err
	}

	if len(resp) == 0 {
		log.Println("add teacher failed", err)
		return 0, errors.New("add teacher failed")
	}

	ret := struct {
		ID int64 `json:"id"`
	}{}

	err = json.Unmarshal(resp, &ret)
	if err != nil {
		log.Println("add teacher failed", err)
		return 0, err
	}
	return ret.ID, nil
}

func modTeacher() {

}

func DeleteAllTeacher(list models.TeacherList) error {
	load := []int64{}
	end := len(list) - 1
	for k, v := range list {
		load = append(load, v.TeacherID)
		if (k+1)%100 == 0 || k == end {
			err := delTeacher(load)
			if err != nil {
				log.Fatal("delete failed")
			}
			load = load[:0]
			log.Println("k", k)
		}
	}
	return nil
}

func delTeacher(id []int64) error {
	resp, err := sendPostRequest(host+"/api/v1/dean/teacher/delete", struct {
		IDList []int64 `json:"id_list"`
	}{IDList: id})
	if err != nil {
		log.Println("delete teacher failed", err)
		return err
	}

	ret := struct {
		FailedList []int64 `json:"failed_list"`
	}{}
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		log.Println("delete failed", err)
	}
	if len(ret.FailedList) > 0 {
		log.Println("failed list", ret.FailedList, "request", id)
	}
	return nil
}
