package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-svr-template/common/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var AuthAes = NewAesCbcPKCS7("XGYUZj78QvlvyHQ1eKeSeNhCJcJRQOyQ")

func GinLogger(threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := log.ColorForStatus(statusCode)
		methodColor := log.ColorForMethod(method)
		userId := c.Request.Header.Get("selfUserId")

		requestData := getRequestData(c)
		log.Infof("[GIN] %s%s%s %s%s %s%d%s %.03f [%s] [user_id:%s] %s",
			methodColor, method, log.Reset,
			c.Request.Host, requestData,
			statusColor, statusCode, log.Reset,
			latency.Seconds(),
			clientIP,
			userId,
			c.Errors.String())

		if latency > threshold {
			log.Warnf("[GIN SLOW] %s%s%s %s%s %s%d%s %.03f [%s] [user_id:%s] startAt: %s endAt: %s",
				methodColor, method, log.Reset,
				c.Request.Host, requestData,
				statusColor, statusCode, log.Reset,
				latency.Seconds(),
				clientIP,
				userId,
				start.Format("15:04:05.999999999"),
				end.Format("15:04:05.999999999"))
		}
	}
}

func getRequestData(c *gin.Context) string {
	var requestData string
	method := c.Request.Method
	if method == "GET" || method == "DELETE" {
		requestData = c.Request.RequestURI
	} else {
		c.Request.ParseForm()
		requestData = fmt.Sprintf("%s [%s]", c.Request.RequestURI, c.Request.Form.Encode())
	}

	if len(requestData) > 1024 {
		return requestData[:1024]
	} else {
		return requestData
	}
}

const (
	STATUS_OK    string = "SUCCESS"
	STATUS_ERROR string = "ERROR"
)

type CommonResHead struct {
	Status  string      `json:"status"`
	Code    int32       `json:"code"`
	Desc    string      `json:"desc"`
	Content interface{} `json:"content"`
}

// status: OK 和 非OK， 非OK的情况下，请把描述填到Desc字段中
// 以后可以转调用SendResponseWithHttpStatus
func SendResponse(c *gin.Context, status string, desc string, data interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	resp := CommonResHead{
		Status:  status,
		Desc:    desc,
		Content: data,
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		log.Error(err.Error())
	} else {
		c.Writer.Write(b)
	}
}

func SendResponseWithHttpStatus(c *gin.Context, status string, desc string, data interface{}, httpStatus int) {
	c.Writer.Header().Set("Content-Type", "application/json")
	resp := CommonResHead{
		Status:  status,
		Desc:    desc,
		Content: data,
	}

	c.JSON(httpStatus, resp)
}

type UserAgent struct {
	AppVersion        string `json:"app_version"`
	MobilePlatform    string `json:"mobile_platform"`
	MobileSystem      string `json:"mobile_system"`
	MobileDeviceBrand string `json:"mobile_device_brand"`
}

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isInPathWhiteList(c.Request.URL.Path) {
			// Process request
			c.Next()
			return
		}

		ua := c.GetHeader("UserAgent")
		if ua == "" {
			log.Warnf("No UserAgent in the req: %+v", c.Request.Header)
			SendResponseWithHttpStatus(c, STATUS_ERROR, "No UserAgent", "", http.StatusUnauthorized)
			c.Abort()
			return
		}

		var userAgent UserAgent
		err := json.Unmarshal([]byte(ua), &userAgent)
		if err != nil {
			log.Errorf("Unmarshal UserAgent Fail: %s %s", ua, err.Error())
			SendResponseWithHttpStatus(c, STATUS_ERROR, "Error UserAgent", "", http.StatusUnauthorized)
			c.Abort()
			return
		}

		authToken := c.GetHeader("Authorization")
		tokenPlaintext, err := AuthAes.Decrypt(authToken)
		if err != nil {
			log.Errorf("Decrypt Authorization error : %s %s", authToken, err.Error())
			SendResponseWithHttpStatus(c, STATUS_ERROR, "Wrong Authorization", "", http.StatusUnauthorized)
			c.Abort()
			return
		}

		items := strings.Split(tokenPlaintext, "|")
		if len(items) != 5 {
			log.Errorf("Wrong1 Authorization: %s", tokenPlaintext)
			SendResponseWithHttpStatus(c, STATUS_ERROR, "Wrong Authorization", "", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if items[0] != userAgent.AppVersion || items[1] != userAgent.MobileSystem || items[2] != userAgent.MobileDeviceBrand {
			log.Errorf("Wrong2 Authorization: %+v", userAgent)
			SendResponseWithHttpStatus(c, STATUS_ERROR, "Wrong Authorization", "", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 判断过期时间
		timeStamp, err := strconv.ParseInt(items[4], 10, 64)
		if err != nil {
			log.Error("wrong time")
			c.Abort()
			return
		}

		tokenTime := time.Unix(timeStamp, 0)
		if tokenTime.Add(time.Duration(time.Hour * 6)).Before(time.Now()) {
			SendResponseWithHttpStatus(c, STATUS_ERROR, "Wrong Authorization", "", http.StatusUnauthorized)
			log.Infof("token timeout. token time: %d", timeStamp)
			c.Abort()
			return
		}

		userId, err := strconv.ParseInt(items[4], 10, 64)
		if err != nil {
			SendResponseWithHttpStatus(c, STATUS_ERROR, "Wrong Authorization", "", http.StatusUnauthorized)
			log.Error("token wrong userId")
			c.Abort()
			return
		}

		err = GenAuth(c, userId)
		if err != nil {
			log.Errorf("gen new auth err: %s", err.Error())
		}

		c.Request.Header.Set("selfUserId", items[3])
	}
}

var mWhitePathMap = map[string]Empty{
	"/wx/wx_login":        empty,
	"/admin/login":        empty,
	"/wechat/wx_callback": empty,
	"/test/ping":          empty,
}

func isInPathWhiteList(path string) bool {
	_, ok := mWhitePathMap[path]
	return ok
}

func GenAuth(c *gin.Context, userId int64) error {
	ua := c.GetHeader("Useragent")
	if ua == "" {
		return fmt.Errorf("No UserAgent")
	}

	var userAgent UserAgent
	err := json.Unmarshal([]byte(ua), &userAgent)
	if err != nil {
		return fmt.Errorf("Unmarshal UserAgent Fail: %s %s", ua, err.Error())
	}

	plaintext := fmt.Sprintf("%s|%s|%s|%d|%d", userAgent.AppVersion, userAgent.MobileSystem, userAgent.MobileDeviceBrand, userId, time.Now().Unix())
	out, err := AuthAes.Encrypt([]byte(plaintext))
	if err != nil {
		return err
	}

	c.Writer.Header().Set("Authorization", out)
	return nil
}

func ShowAllRawBody(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	log.Infof("Raw Http Body: %s", string(data))
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Useragent,selfuserid")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin,"+
			" Access-Control-Allow-Headers, Content-Type,Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
