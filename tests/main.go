package main

import (
	"log"
)

var host = "http://127.0.0.1:2008"

func main() {
	// 1. login
	_, err := login(host, "aronic", "123456")
	if err != nil {
		log.Fatal(err)
	}

	// remove all data
	subject, err := GetAllSubjects()
	log.Println(len(subject))

	// list
	list, err := GetTeacherListAll(host)
	if err != nil {
		return
	}

	log.Println("total teacher num:", len(list))

	err = DeleteAllSubject(subject)
	err = DeleteAllTeacher(list)

	err = TeacherNormalFlow()
	if err != nil {
		log.Println(err)
		return
	}

	err = TeacherAbnormalFlow()
	if err != nil {
		log.Println(err)
		return
	}
	subject = AddMultiSubject()

	err = AddMultiTeachers(subject)
	if err != nil {
		log.Println("add teacher failed", err)
		return
	}

}
