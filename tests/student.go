package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/arong/dean/base"

	"github.com/arong/dean/models"
	"github.com/pkg/errors"
)

func ImportStudent(host string) {
	for i := 0; i < 10000; i++ {
		s := getFakeStudent(host)
		id, err := addStudent(s)
		if err != nil || id == 0 {
			log.Println("add student failed", "err", err)

		}
	}
}

func GetAllStudent(host string) models.StudentList {
	i := 1
	all := models.StudentList{}
	for {
		list, err := getStudentList(host, i, 100)
		if err != nil {
			break
		}
		all = append(all, list...)
		if len(list) < 100 {
			break
		}
		i++
	}
	return all
}

func getFakeStudent(host string) models.StudentInfo {
	resp, err := http.Get("http://192.168.231.132:8080/api")
	if err != nil {
		log.Println("random source failed, check it")
		return models.StudentInfo{}
	}

	if err != nil {
		return models.StudentInfo{}
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.StudentInfo{}
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
		return models.StudentInfo{}
	}
	return models.StudentInfo{
		Name:       ret.Name,
		Gender:     rand.Intn(2) + 1,
		Mobile:     ret.Mobile,
		Address:    ret.Addr,
		RegisterID: getRegisterID(),
	}
}

func addStudent(student models.StudentInfo) (int64, error) {
	resp, err := sendPostRequest(host+"/api/v1/dean/student/add", student)
	if err != nil {
		log.Println("add student failed", err)
		return 0, err
	}

	if len(resp) == 0 {
		log.Println("add student failed", err)
		return 0, errors.New("add student failed")
	}

	ret := struct {
		ID int64 `json:"id"`
	}{}

	err = json.Unmarshal(resp, &ret)
	if err != nil {
		log.Println("add student failed", err)
		return 0, err
	}
	return ret.ID, nil
}

func getStudentList(host string, page, size int) (models.StudentList, error) {
	resp, err := sendPostRequest(host+"/api/v1/dean/student/filter", base.CommPage{Page: page, Size: size})
	if err != nil {
		log.Println("sendPostRequest failed", err)
		return nil, err
	}

	ret := struct {
		Total int                `json:"total"`
		List  models.StudentList `json:"list"`
	}{}

	//fmt.Println(string(resp))
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		log.Println("unrecognized return value", err)
		return nil, err
	}

	return ret.List, nil
}

func getRegisterID() string {
	year := 2017 + rand.Intn(2)
	class := rand.Intn(99) + 1
	number := rand.Intn(100)

	registerID := fmt.Sprintf("%d%02d%02d", year, class, number)
	return registerID
}

func deleteAllStudent(list models.StudentList) error {
	load := []int64{}
	end := len(list) - 1
	for k, v := range list {
		load = append(load, v.StudentID)
		if (k+1)%100 == 0 || k == end {
			err := deleteStudent(load)
			if err != nil {
				log.Fatal("delete failed")
			}
			load = load[:0]
			log.Println("k", k)
		}
	}
	return nil
}

func deleteStudent(id []int64) error {
	resp, err := sendPostRequest(host+"/api/v1/dean/student/delete", struct {
		IDList []int64 `json:"id_list"`
	}{IDList: id})
	if err != nil {
		log.Println("delete student failed", err)
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

func getStudentInfo(id int64) (models.StudentInfo, error) {
	var ret models.StudentInfo
	resp, err := sendGetRequest(host+"/api/v1/dean/student/info/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
