package apis

import (
	"encoding/json"
	"fmt"
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/chanxuehong/wechat/mp/message/callback/request"
	"github.com/chanxuehong/wechat/mp/message/callback/response"
	"github.com/gin-gonic/gin"
	"github.com/diffguo/gocom/log"
	"go-svr-template/models"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var mutex sync.Mutex

var (
	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	UserMsgHandler        core.Handler
	UserMsgServer         *core.Server
	UserClient            *core.Client
	UserAccessTokenServer core.AccessTokenServer
	UserWeChatClient      *core.Client
)

type GetWxTokenRsp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type WeChatPublicServer struct {
	AppId     string
	AppSecret string
	core.AccessTokenServer
}

func (r WeChatPublicServer) Token() (string, error) {
	rsp := r.GetWxTokenDo("")
	if rsp.ErrMsg != "" {
		return "", fmt.Errorf("%s", rsp.ErrMsg)
	}
	return rsp.AccessToken, nil
}

func (r WeChatPublicServer) RefreshToken(currentToken string) (string, error) {
	rsp := r.GetWxTokenDo(currentToken)
	if rsp.ErrMsg != "" {
		return "", fmt.Errorf("%s", rsp.ErrMsg)
	}
	return rsp.AccessToken, nil
}

func (WeChatPublicServer) IID01332E16DF5011E5A9D5A4DB30FED8E1() {}

func (r WeChatPublicServer) GetWxTokenDo(oldToken string) (rsp GetWxTokenRsp) {
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	token := models.TWeChatAccessToken{}
	err = token.FirstOrCreate(nil, r.AppId)
	if err == nil {
		if token.ExpireAT.Unix() > time.Now().Unix() {
			if oldToken == "" || oldToken != token.AccessToken {
				rsp.AccessToken = token.AccessToken
				return
			}
		}
	} else {
		log.Errorf(err.Error())
	}

	reqURL := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + r.AppId + "&secret=" + r.AppSecret
	for {
		resp, err := http.Get(reqURL)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		resp.Body.Close()

		if err == nil {
			log.Infof("refresh %s %s oldToken %s", reqURL, data, oldToken)
			err = json.Unmarshal(data, &rsp)
			if err == nil {
				if rsp.ErrCode == 0 {
					rsp.ExpiresIn = rsp.ExpiresIn - 5*60
					token.Update(nil, rsp.AccessToken, time.Unix(time.Now().Unix()+int64(rsp.ExpiresIn), 0))
					return
				}

				log.Error(rsp.ErrMsg)
			} else {
				log.Error(err.Error())
			}
		}
		time.Sleep(time.Second)

	}

	return
}

func InitUserWeChat() {
	return

	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(UserDefaultMsgHandler)
	mux.DefaultEventHandleFunc(UserDefaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, UserTextMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, UserMenuClickEventHandler)

	UserMsgHandler = mux

	UserMsgServer = core.NewServer("gh_d44", "wx3a376", "tow",
		"ssov6ZXfeiPtukQTfTQssOea1lMuqlcAZagp8KSivss", UserMsgHandler, nil)

	UserAccessTokenServer = WeChatPublicServer{AppId: "wx3a3762ddddd53975", AppSecret: "f6c387977fb63f23e4c11111a9c08d02"}
	UserWeChatClient = core.NewClient(UserAccessTokenServer, nil)

	token, err := UserWeChatClient.AccessTokenServer.Token()
	if err != nil {
		log.Errorf("获取user token失败:%s ", err.Error())
	} else {
		log.Infof("获取user token成功: %s", token)
	}
}

func UserTextMsgHandler(ctx *core.Context) {
	log.Infof("收到文本消息:\n%s\n", ctx.MsgPlaintext)
	msg := request.GetText(ctx.MixedMsg)

	if len(msg.Content) == 11 {
		log.Infof("收到电话：%s", msg.Content)
	}

	resp := response.NewText(ctx.MixedMsg.FromUserName, ctx.MixedMsg.ToUserName, int64(time.Now().Unix()), "回复预订电话号码，查看您的信息")
	ctx.RawResponse(resp)
}

func UserDefaultMsgHandler(ctx *core.Context) {
	log.Infof("收到消息:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func UserMenuClickEventHandler(ctx *core.Context) {
	log.Infof("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "回复预订电话号码，查看您的信息")
	ctx.RawResponse(resp) // 明文回复
}

func UserDefaultEventHandler(ctx *core.Context) {
	log.Infof("收到事件:\n%s\n", ctx.MsgPlaintext)

	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "回复预订电话号码，查看您的信息")
	err := ctx.RawResponse(resp) // 明文回复
	if err != nil {
		log.Infof("%s", err.Error())
	}
}

func UserWxCallbackHandler(c *gin.Context) {
	UserMsgServer.ServeHTTP(c.Writer, c.Request, nil)
}
