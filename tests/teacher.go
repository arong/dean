package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

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

func AddMultiTeachers(list models.SubjectList) error {
	size := len(list)

	// get basic info
	for i := 0; i < 10000; i++ {
		t := getFakeInfo()
		t.SubjectID = list[rand.Intn(size)].ID
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
	resp, err := http.Get("http://192.168.231.132:8080/api")
	if err != nil {
		log.Println("random source failed, check it")
		return models.Teacher{}
	}

	if err != nil {
		return models.Teacher{}
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
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

	err = json.Unmarshal(buff, &ret)
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

func modTeacher(teacher models.Teacher) error {
	_, err := sendPostRequest(host+"/api/v1/dean/teacher/modify", teacher)
	if err != nil {
		log.Println("add teacher failed", err)
		return err
	}
	return nil
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

	if len(resp) > 0 {
		err = json.Unmarshal(resp, &ret)
		if err != nil {
			log.Println("delete failed", err)
		}
		if len(ret.FailedList) > 0 {
			log.Println("failed list", ret.FailedList, "request", id)
		}
	}
	return nil
}

func TeacherNormalFlow() error {
	t := models.Teacher{
		TeacherMeta: models.TeacherMeta{Name: "无名", Gender: 1, Birthday: "1998-09-06"},
	}
	id, err := addTeacher(t)
	if err != nil {
		return err
	}
	log.Println("test add success")

	ret, err := getTeacherInfo(id)
	if err != nil {
		return err
	}

	if !ret.Equal(t.TeacherMeta) {
		log.Println("info not match, ret", ret, "t", t)
		return errors.New("info not match")
	}
	log.Println("test get info success")

	// modify
	ret.Address = strings.Repeat("烫", 64)
	ret.Birthday = "1987-09-06"
	ret.Mobile = "18765241623"

	err = modTeacher(*ret)
	if err != nil {
		log.Println("modify teacher failed", err)
		return err
	}

	curr, err := getTeacherInfo(ret.TeacherID)
	if err != nil {
		log.Println("failed to query", err)
		return err
	}

	if !curr.Equal(ret.TeacherMeta) {
		log.Println("info not match, ret", ret, "t", t)
		return errors.New("info not match")
	}
	log.Println("test modify success")

	err = delTeacher([]int64{id})
	if err != nil {
		return err
	}
	log.Println("test delete success")

	_, err = getTeacherInfo(id)
	if err == nil {
		log.Println("logic failure", err)
		return err
	}
	log.Println("test get success")

	list, err := GetTeacherListAll(host)
	if err != nil {
		log.Println("get list failed", err)
		return err
	}

	if len(list) != 0 {
		log.Println("list logic failed")
	}
	log.Println("test list success")
	return nil
}

func TeacherAbnormalFlow() error {
	t := models.Teacher{
		TeacherMeta: models.TeacherMeta{Name: "无名", Gender: 1, Birthday: "1998-09-06", SubjectID: 1999},
	}

	_, err := addTeacher(t)
	if err == nil {
		return errors.New("logic error")
	}

	return nil
}
