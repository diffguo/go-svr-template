package apis

import (
	"github.com/diffguo/gocom"
	"github.com/diffguo/gocom/log"
	"github.com/gin-gonic/gin"
	"go-svr-template/io"
	"go-svr-template/models"
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
func PingPong(c *gin.Context) {
	type InputStructure struct {
		Content     string `form:"content" binding:"required,len=2"`
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

//https://godoc.org/gopkg.in/go-playground/validator.v9
func TestValidator(c *gin.Context) {
	type InputStructure struct {
		Content     string `form:"content" binding:"required,len=2"`
		Age         int    `form:"age" binding:"max=10,min=1"`
		Pass        string `form:"pass" binding:"gte=6"`           // 字母和数字
		ConfirmPass string `form:"cpass" binding:"eqfield=Pass"`   // 与上一个域相同
		FileName    string `form:"fname" binding:"gte=2,alphanum"` // 字母和数字
		Email       string `form:"email" binding:"email"`
		//Body    string `form:"base64" binding:"base64"` // startswith=hello, endswith=hello, contains=@, uuid, ip
	}

	var is InputStructure
	ok := gocom.Bind(c, &is)
	if !ok {
		io.SendResponse(c, "", io.ErrCodeParamErr)
		return
	}

	log.Infof("TestValidator: %+v", is)
	gocom.SendSimpleResponse(c, is.Content)
}

func TransactionDemo(c *gin.Context) {
	var err error
	tx := models.LocalDB{DB: models.GDB.Begin()}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	comment := models.TComment{Content: "t1"}
	user := models.User{Name: "testName"}

	err = comment.Create(&tx)
	if err != nil {
		log.Errorf("comment.Create err: %s", err.Error())
		gocom.SendSimpleResponse(c, "TransactionDemo err")
		return
	}

	err = user.Create(&tx)
	if err != nil {
		log.Errorf("user.Create err: %s", err.Error())
		gocom.SendSimpleResponse(c, "TransactionDemo err")
		return
	}

	gocom.SendSimpleResponse(c, "TransactionDemo Success")
}
