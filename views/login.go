package views

import (
	"github.com/gin-gonic/gin"
	"go-svr-template/common"
	"go-svr-template/common/log"
	"go-svr-template/controller"
	"go-svr-template/models"
)

func ApiCheckAuth(c *gin.Context) {
	_, err := controller.GetSelfUserId(c)
	if err != nil {
		return
	} else {
		common.SendResponse(c, common.STATUS_OK, "", "")
	}
}

func ApiLogin(c *gin.Context) {
	type InputStructure struct {
		MobileNumber string `json:"mobile_number"`
		PassWord     string `json:"pass_word"`
	}

	var is InputStructure
	err := c.Bind(&is)
	if err != nil {
		log.Errorf("miss code")
		common.SendResponse(c, common.STATUS_ERROR, "miss code", "")
		return
	}

	passWord := controller.Hmac4Password(is.PassWord)
	up, err := models.GetUserByPassword(nil, is.MobileNumber, passWord)
	if err != nil {
		common.SendResponse(c, common.STATUS_ERROR, err.Error(), "")
		return
	}

	err = common.GenAuth(c, up.ID)
	if err != nil {
		common.SendResponse(c, common.STATUS_ERROR, err.Error(), "")
		return
	}
	common.SendResponse(c, common.STATUS_OK, "", up)
}
