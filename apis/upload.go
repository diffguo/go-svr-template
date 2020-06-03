package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-svr-template/controller"
	"go-svr-template/common"
	"go-svr-template/common/log"
	"go-svr-template/common/tools"
)

const CDNUrl = "https://xxx.oss-cn-chengdu.aliyuncs.com"

func ApiUploadAvatar(c *gin.Context) {
	userId, err := controller.GetSelfUserId(c)
	if err != nil {
		return
	}

	postFix := c.Param("post_fix")
	resourcePath := fmt.Sprintf("user/avatar/%d.%s", userId, postFix)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Errorf("read body err: %s", err.Error())
		common.SendResponse(c, common.STATUS_ERROR, "上传失败，读取上传数据失败!", "")
		return
	}

	contentType := header.Header["Content-Type"][0]
	log.Infof("upload file: %s size: %d, userId: %d, to file: %s", contentType, header.Size, userId, resourcePath)

	success := tools.UploadToTWNoExpireOss(resourcePath, contentType, file)
	if success {
		common.SendResponse(c, common.STATUS_OK, "", fmt.Sprintf("%s/%s", CDNUrl, resourcePath))
	} else {
		common.SendResponse(c, common.STATUS_ERROR, "上传失败，请重试", "")
	}
}
