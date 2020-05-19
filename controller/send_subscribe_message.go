package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SendSubscribeMessageReqDataItem struct {
	Value string `json:"value"`
}

type SendSubscribeMessageReq struct {
	Data       map[string]SendSubscribeMessageReqDataItem `json:"data"`
	Page       string                                     `json:"page"`
	TemplateID string                                     `json:"template_id"`
	Touser     string                                     `json:"touser"`
}

//发送订阅消息
func SendSubscribeMessage(accessToken string, req SendSubscribeMessageReq) (body string, err error) {
	bytesData, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	url := "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=" + accessToken
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	body = string(respBytes)
	return
}
