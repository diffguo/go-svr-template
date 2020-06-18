package apis

import (
	"fmt"
	"github.com/diffguo/gocom"
	"github.com/diffguo/gocom/log"
	"github.com/diffguo/gocom/tools"
	"github.com/gin-gonic/gin"
	"go-svr-template/io"
)

const CDNUrl = ""
var YourBucket *tools.OssBucket

func InitOss()  {
	var err error
	YourBucket, err = tools.InitOssBucket("", "", "", "", 100)
	if err != nil {
		fmt.Printf("init oss err: %s", err.Error())
	}
}

func ApiUploadAvatar(c *gin.Context) {
	userId, err := io.GetSelfUserId(c)
	if err != nil {
		return
	}

	postFix := c.Param("post_fix")
	resourcePath := fmt.Sprintf("user/avatar/%d.%s", userId, postFix)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Errorf("read body err: %s", err.Error())
		gocom.SendResponseImp(c, "", io.ErrCodeParamErr, "上传失败，读取上传数据失败!")
		return
	}

	contentType := header.Header["Content-Type"][0]
	log.Infof("upload file: %s size: %d, userId: %d, to file: %s", contentType, header.Size, userId, resourcePath)

	success := YourBucket.UploadToOss(resourcePath, contentType, file)
	if success {
		gocom.SendSimpleResponse(c, fmt.Sprintf("%s/%s", CDNUrl, resourcePath))
	} else {
		gocom.SendResponseImp(c, "", io.ErrCodeLogicErr, "上传到OSS失败，请重试")
	}
}
