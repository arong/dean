package main

import (
	"fmt"
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
	// remove subject
	subject, err := GetAllSubjects()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("total subject", len(subject))
	err = DeleteAllSubject(subject)
	if err != nil {
		log.Println(err)
		return
	}

	//remove teacher
	list, err := GetTeacherListAll(host)
	if err != nil {
		return
	}

	log.Println("total teacher num:", len(list))

	//err = DeleteAllTeacher(list)

	//err = TeacherNormalFlow()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	//err = TeacherAbnormalFlow()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	// add data back
	subject = AddMultiSubject()
	for _, v := range subject {
		log.Println(v.ID)
	}
	subject, err = GetAllSubjects()
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range subject {
		log.Println(v)
	}
	err = AddMultiTeachers(subject)
	if err != nil {
		log.Println("add teacher failed", err)
		return
	}

}

func testStudent() {
	// 1. login
	_, err := login(host, "aronic", "123456")
	if err != nil {
		log.Fatal(err)
	}
	students := GetAllStudent(host)
	fmt.Println("student number is ", len(students))

	for _, v := range students {
		tmp, err := getStudentInfo(v.StudentID)
		if err != nil {
			log.Fatal("[] logic error, not exist")
		}
		if !tmp.Equal(v) {
			log.Fatal("[] logic error, not equal")
		}
	}

	err = deleteAllStudent(students)
	if err != nil {
		log.Println("deleteAllStudent failed", err)
	}

	ImportStudent(host)
	students = GetAllStudent(host)
	fmt.Println("student number is ", len(students))
}
