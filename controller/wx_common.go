package controller

import (
	"encoding/json"
	"fmt"
	"github.com/diffguo/gocom/log"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	WXAppId            = ""
	WXAppSecret        = ""
)

var accessToken string
var expiresIn float64 = 10
var lastTokenTime time.Time

func GetWXAccessToken() string {
	if time.Now().Sub(lastTokenTime).Seconds() >= expiresIn {
		log.Debug("get wx access token from remote")
		lastTokenTime = time.Now()

		accessToken, expiresIn = _innerGetWXAccessToken()
		log.Debugf("WX Access token: %s %f", accessToken, expiresIn)
	}

	return accessToken
}

func _innerGetWXAccessToken() (string, float64) {
	resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		WXAppId, WXAppSecret))
	if err != nil {
		log.Errorf("get wx token error: %s", err.Error())
		return "", 0
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("get wx token error: %s", err.Error())
		return "", 0
	}

	log.Infof("get wx token: %s", string(body))

	type WXToken struct {
		Access_token string
		Expires_in   int
	}

	ret := WXToken{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		log.Errorf("get wx token error: %s", err.Error())
		return "", 0
	}

	return ret.Access_token, float64(ret.Expires_in)
}
