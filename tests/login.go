package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
)

func login(host, name, password string) (string, error) {
	url := "/api/v1/auth/login"
	data := struct {
		LoginName string `json:"login_name"`
		Password  string `json:"password"`
		Type      int    `json:"type"`
	}{
		LoginName: name,
		Password:  password,
		Type:      2,
	}

	data.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	buff, err := sendPostRequest(host+url, data)
	if err != nil {
		log.Println("sendPostRequest failed", err)
		return "", err
	}

	// set global token
	ret := struct {
		Token string `json:"token"`
	}{}
	err = json.Unmarshal([]byte(buff), &ret)
	if err != nil {
		log.Println("unrecognized return value", string(buff))
		return "", err
	}
	token = ret.Token

	return token, nil
}
