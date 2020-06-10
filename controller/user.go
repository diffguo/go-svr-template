package controller

import (
	"github.com/gin-gonic/gin"
	"go-svr-template/common"
	"go-svr-template/common/log"
	"go-svr-template/models"
)

func GetUser(c *gin.Context) *models.User {
	userId, err := GetSelfUserId(c)
	if err != nil {
		return nil
	}

	user, err := models.GetUserByUserId(nil, userId)
	if err != nil {
		log.Error("Get User By Id err: %s, userId: %d", err.Error(), userId)
		common.SendResponseImp(c, "", common.ErrCodeDBErr, "GetUserByUserId err")
		return nil
	}

	return user
}

