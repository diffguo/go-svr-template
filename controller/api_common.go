package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-svr-template/common"
	"go-svr-template/common/log"
	"strconv"
)

type ObjProfileContainer struct {
	ActionType string
	StrValue   string
	ExtValue   string
}

func GetUserIdFromHeader(c *gin.Context) (int64, error) {
	userId := c.Request.Header.Get("selfUserId")
	if userId != "" {
		iUserId, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("Wrong UserId In Header: %s", userId)
		} else {
			return iUserId, nil
		}
	}

	return 0, fmt.Errorf("no userId")
}

func GetSelfUserId(c *gin.Context) (int64, error) {
	var err error
	var userId int64
	defer func() {
		if err != nil || userId == 0 {
			log.Errorf("cant get userId: %s", err.Error())
			common.SendResponseImp(c, "", common.ErrCodeParamErr, "Get User Id err")
		}
	}()

	userId, err = GetUserIdFromHeader(c)
	if err == nil {
		log.Debugf("get userId from header, userId: %d", userId)
		return userId, nil
	} else {
		return 0, err
	}
}
