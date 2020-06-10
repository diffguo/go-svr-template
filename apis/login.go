package apis

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
		common.SendResponse(c, "", common.ErrCodeParamErr)
	}
}

func ApiLogin(c *gin.Context) {
	type InputStructure struct {
		MobileNumber string `json:"mobile_number"`
		PassWord     string `json:"pass_word"`
	}

	var is InputStructure
	ok := common.Bind(c, &is)
	if !ok {
		common.SendResponse(c, "", common.ErrCodeParamErr)
		return
	}

	passWord := controller.Hmac4Password(is.PassWord)
	up, err := models.GetUserByPassword(nil, is.MobileNumber, passWord)
	if err != nil {
		log.Errorf("db err: %s", err.Error())
		common.SendResponse(c, "", common.ErrCodeDBErr)
		return
	}

	err = common.GenAuth(c, up.ID)
	if err != nil {
		common.SendResponseImp(c, "", common.ErrCodeLogicErr, "GenAuth Error")
		return
	}

	common.SendSimpleResponse(c, up)
}
