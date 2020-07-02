package controller

import (
	"github.com/diffguo/gocom"
	"github.com/diffguo/gocom/log"
	"github.com/gin-gonic/gin"
	"go-svr-template/io"
	"go-svr-template/models"
)

func GetUser(c *gin.Context) *models.TUser {
	userId, err := io.GetSelfUserId(c)
	if err != nil {
		return nil
	}

	user := models.TUser{ID: userId}
	err = models.FindFirst(nil, &user, "ID")
	if err != nil {
		log.Error("Get User By Id err: %s, userId: %d", err.Error(), userId)
		gocom.SendResponseImp(c, "", io.ErrCodeDBErr, "GetUserByUserId err")
		return nil
	}

	return &user
}

