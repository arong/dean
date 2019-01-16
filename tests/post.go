package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var token string

func sendPostRequest(URL string, data interface{}) (json.RawMessage, error) {
	req := struct {
		Token     string      `json:"token"`
		Timestamp int64       `json:"timestamp"`
		Data      interface{} `json:"data"`
		Check     string      `json:"check"`
	}{
		Token:     token,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
	req.Check = fmt.Sprintf("%x", md5.Sum([]byte(req.Token+strconv.FormatInt(req.Timestamp, 10))))
	buff, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(buff))
	request, err := http.NewRequest("POST", URL, strings.NewReader(string(buff)))
	if err != nil {
		return nil, err
	}
	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	buff, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}{}

	err = json.Unmarshal(buff, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, fmt.Errorf("%s", ret.Msg)
	}

	return ret.Data, nil
}

func sendGetRequest(URL string, data interface{}) ([]byte, error) {
	req := struct {
		Token     string `json:"token"`
		Timestamp int64  `json:"timestamp"`
		Check     string `json:"check"`
	}{
		Token:     token,
		Timestamp: time.Now().Unix(),
	}
	req.Check = fmt.Sprintf("%x", md5.Sum([]byte(req.Token+strconv.FormatInt(req.Timestamp, 10))))

	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	q.Set("token", req.Token)
	q.Set("timestamp", strconv.FormatInt(req.Timestamp, 10))
	q.Set("check", req.Check)
	request.URL.RawQuery = q.Encode()

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{}

	err = json.Unmarshal(buff, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, fmt.Errorf("%s", ret.Msg)
	}

	buff, err = json.Marshal(ret.Data)
	if err != nil {
		return nil, err
	}
	return buff, nil
}
