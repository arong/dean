package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
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
	fmt.Println(data.Password)

	buff, err := sendPostRequest(host+url, data)
	if err != nil {
		return "", err
	}

	// set global token
	token = string(buff)
	token = strings.Replace(token, "\"", "", -1)

	return token, nil
}
