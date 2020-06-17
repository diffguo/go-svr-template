package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/diffguo/gocom"
	"github.com/diffguo/gocom/log"
	"go-svr-template/io"
)

// PingPong godoc
// @Summary Test Sever is working
// @Description http://127.0.0.1:8010/test/ping?content=1111
// @Accept  json
// @Produce  json
// @Param content path string false "ping pong content"
// @Success 200 {object} common.CommonRspHead
// @Failure 400 {object} common.CommonRspHead
// @Failure 404 {object} common.CommonRspHead
// @Failure 500 {object} common.CommonRspHead
// @Router /test/ping [get]
func PingPong(c *gin.Context)  {
	type InputStructure struct {
		Content string `form:"content" binding:"required"`
	}

	var is InputStructure
	ok := gocom.Bind(c, &is)
	if !ok {
		io.SendResponse(c, "", io.ErrCodeParamErr)
		return
	}

	log.Infof("PingPong: %+v", is)
	gocom.SendSimpleResponse(c, is.Content)
}