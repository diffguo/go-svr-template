package io

import (
	"github.com/diffguo/gocom"
	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, content interface{}, errCode int) {
	if errCode == 0 {
		gocom.SendResponseImp(c, content, 0, "")
		return
	}

	if errCode >= ErrCodeParamErr {
		gocom.SendResponseImp(c, content, errCode, MapErrCode2Desc[errCode])
		return
	}

	gocom.SendResponseImp(c, content, errCode, "")
}