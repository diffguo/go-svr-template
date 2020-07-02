package apis

import (
	"github.com/diffguo/gocom"
	"github.com/diffguo/gocom/log"
	"github.com/diffguo/gocom/tools"
	"github.com/diffguo/gocom/wx_pay"
	"github.com/gin-gonic/gin"
	"go-svr-template/controller"
	"go-svr-template/io"
	"go-svr-template/models"
	"net/http"
)

func ApiDecodePhoneNumber(c *gin.Context) {
	userId, err := io.GetSelfUserId(c)
	if err != nil {
		return
	}

	userWX := models.TUserWX{UserId: userId}
	err = models.FindFirst(nil, &userWX, "UserId")
	if err != nil {
		log.Errorf("cant get user profile: %d", userId)
		io.SendResponse(c, "", io.ErrCodeDBErr)
		return
	}

	type InputStructure struct {
		EncryptedData string `json:"encrypted_data"`
		Iv            string `json:"iv"`
	}

	ts := InputStructure{}
	err = c.Bind(&ts)
	if err != nil || ts.EncryptedData == "" || ts.Iv == "" {
		log.Errorf("miss param: %+v, userId: %d", ts, userId)
		io.SendResponse(c, "", io.ErrCodeParamErr)
		return
	}

	log.Infof("ApiDecodePhoneNumber: %+v", ts)

	pc := tools.WxBizDataCrypt{AppID: controller.WXAppId, SessionKey: userWX.WXSessionKey}
	result, err := pc.Decrypt(ts.EncryptedData, ts.Iv, false)
	if err == nil {
		mobileNumber := result.(map[string]interface{})["purePhoneNumber"].(string)

		user := models.TUser{ID: userId}
		paras := map[string]interface{}{"MobileNumber": mobileNumber}
		if err = models.Update(nil, &user, paras, "ID"); err != nil {
			log.Errorf("UpdateUserByUserId Err: %s, userId: %d", err.Error(), userId)
			io.SendResponse(c, "", io.ErrCodeDBErr)
			return
		} else {
			gocom.SendSimpleResponse(c, mobileNumber)
		}
	} else {
		log.Errorf("PhoneNumber decode Err: %s, userId: %d", err.Error(), userId)
		gocom.SendResponseImp(c, "", io.ErrCodeLogicErr, "PhoneNumber decode Err")
	}
}

func ApiWXPayCallBack(c *gin.Context) {
	log.Info("收到微信支付回调")
	wxPayParams := wx_pay.DecodeWXPayParamsFromXML(c.Request.Body)

	log.Info("参数：", wxPayParams)

	okContent := "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
	failContent := "<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[签名失败]]></return_msg></xml>"

	if wxPayParams.GetString("result_code") != "SUCCESS" {
		log.Warnf("支付失败. 参数：%+v", wxPayParams)
		c.String(http.StatusOK, okContent)
		return
	}

	sign := wx_pay.WXPayClient.Sign(wxPayParams)
	if sign != wxPayParams.GetString("sign") {
		log.Warnf("验签失败. 参数：%+v, sign: %s", wxPayParams, sign)
		c.String(http.StatusOK, failContent)
		return
	}

	c.String(http.StatusOK, okContent)
	return
}
