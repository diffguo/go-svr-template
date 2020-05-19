package views

import (
	"github.com/gin-gonic/gin"
	"go-svr-template/common"
	"go-svr-template/common/log"
)

// http://127.0.0.1:8010/test/ping?content=1111
func PingPong(c *gin.Context)  {
	type InputStructure struct {
		Content string `form:"content"`
	}

	ts := InputStructure{}
	err := c.Bind(&ts)
	if err != nil {
		log.Errorf("bind err: %d", err.Error())
		common.SendResponse(c, common.STATUS_ERROR, "bind err", "")
		return
	}

	log.Infof("PingPong: %+v", ts)
	common.SendResponse(c, common.STATUS_OK, "pong", ts.Content)
}