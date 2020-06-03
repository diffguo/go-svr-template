package apis

import (
	"github.com/gin-gonic/gin"
	"go-svr-template/common"
	"go-svr-template/common/log"
)

// PingPong godoc
// @Summary Test Sever is working
// @Description http://127.0.0.1:8010/test/ping?content=1111
// @Accept  json
// @Produce  json
// @Param content path string false "ping pong content"
// @Success 200 {object} common.CommonResHead
// @Failure 400 {object} common.CommonResHead
// @Failure 404 {object} common.CommonResHead
// @Failure 500 {object} common.CommonResHead
// @Router /test/ping [get]
func PingPong(c *gin.Context)  {
	type InputStructure struct {
		Content string `form:"content" binding:"required"`
	}

	var is InputStructure
	ok := common.Bind(c, &is)
	if !ok {
		common.SendResponse(c, common.STATUS_ERROR, "param err", "")
		return
	}

	log.Infof("PingPong: %+v", is)
	common.SendResponse(c, common.STATUS_OK, "pong", is.Content)
}