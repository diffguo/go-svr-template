package apis

import (
	"github.com/diffguo/gocom"
	"github.com/diffguo/gocom/log"
	"github.com/gin-gonic/gin"
	"go-svr-template/controller"
	"go-svr-template/io"
	"go-svr-template/models"
)

func ApiCheckAuth(c *gin.Context) {
	_, err := io.GetSelfUserId(c)
	if err != nil {
		return
	} else {
		io.SendResponse(c, "", io.ErrCodeParamErr)
	}
}

func ApiLogin(c *gin.Context) {
	type InputStructure struct {
		MobileNumber string `json:"mobile_number"`
		Password     string `json:"password"`
	}

	var is InputStructure
	ok := gocom.Bind(c, &is)
	if !ok {
		io.SendResponse(c, "", io.ErrCodeParamErr)
		return
	}

	passWord := controller.Hmac4Password(is.Password)
	user := models.TUser{MobileNumber: is.MobileNumber, Password: passWord}
	err := models.FindFirst(nil, &user, "MobileNumber", "Password")
	if err != nil {
		log.Errorf("db err: %s", err.Error())
		io.SendResponse(c, "", io.ErrCodeDBErr)
		return
	}

	err = gocom.GenAuth(c, user.ID)
	if err != nil {
		gocom.SendResponseImp(c, "", io.ErrCodeLogicErr, "GenAuth Error")
		return
	}

	gocom.SendSimpleResponse(c, user)
}
