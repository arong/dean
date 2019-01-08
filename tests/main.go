package main

import (
	"fmt"
	"log"
)

var host = "http://127.0.0.1:2008"

func main() {
	// 1. login
	token, err := login(host, "aronic", "123456")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	// list
	list, err := getTeacherList(host)
	if err != nil {
		return
	}

	for _, v := range list {
		info, err := getTeacherInfo(v.TeacherID)
		if err != nil {
			fmt.Println(err)
		}
		if info.Name != v.Name || info.TeacherID != v.TeacherID || info.SubjectID != v.SubjectID {
			fmt.Println("failed")
			break
		}
	}

}
