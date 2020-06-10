package apis

import (
	"github.com/gin-gonic/gin"
	"go-svr-template/common"
	"go-svr-template/common/log"
	"go-svr-template/common/tools"
	"go-svr-template/common/wx_pay"
	"go-svr-template/controller"
	"go-svr-template/models"
	"net/http"
	"time"
)

func ApiDecodePhoneNumber(c *gin.Context) {
	userId, err := controller.GetSelfUserId(c)
	if err != nil {
		return
	}

	upWX, err := models.GetUserWXWithUserId(nil, userId)
	if err != nil {
		log.Errorf("cant get user profile: %d", userId)
		common.SendResponse(c, "", common.ErrCodeDBErr)
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
		common.SendResponse(c, "", common.ErrCodeParamErr)
		return
	}

	log.Infof("ApiDecodePhoneNumber: %+v", ts)

	pc := tools.WxBizDataCrypt{AppID: controller.WXAppId, SessionKey: upWX.WXSessionKey}
	result, err := pc.Decrypt(ts.EncryptedData, ts.Iv, false)
	if err == nil {
		mobileNumber := result.(map[string]interface{})["purePhoneNumber"].(string)
		if err = models.UpdateUserByUserId(nil, userId, map[string]interface{}{"mobile_number": mobileNumber, "mobile_verified": 1, "mobile_verify_time": time.Now()}); err != nil {
			log.Errorf("UpdateUserByUserId Err: %s, userId: %d", err.Error(), userId)
			common.SendResponse(c, "", common.ErrCodeDBErr)
			return
		} else {
			common.SendSimpleResponse(c, mobileNumber)
		}
	} else {
		log.Errorf("PhoneNumber decode Err: %s, userId: %d", err.Error(), userId)
		common.SendResponseImp(c, "", common.ErrCodeLogicErr, "PhoneNumber decode Err")
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

	client := wx_pay.Client{ApiKey: wx_pay.WX_APP_KEY}
	sign := client.Sign(wxPayParams)
	if sign != wxPayParams.GetString("sign") {
		log.Warnf("验签失败. 参数：%+v, sign: %s", wxPayParams, sign)
		c.String(http.StatusOK, failContent)
		return
	}

	c.String(http.StatusOK, okContent)
	return
}
