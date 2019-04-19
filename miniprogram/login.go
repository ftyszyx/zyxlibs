package miniprogram

import (
	"encoding/json"
	"fmt"

	"github.com/ftyszyx/libs/wechat/util"
	"github.com/pkg/errors"
)

const (
	getsession_url = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type ResAccessToken struct {
	Openid      string `json:"openid"`
	Session_key string `json:"session_key"`
	Unionid     string `json:"unionid"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

// GetUserAccessToken 通过网页授权的code 换取access_token(区别于context中的access_token)
func (oauth *MiniData) Getcode2Session(code string) (result ResAccessToken, err error) {
	urlStr := fmt.Sprintf(getsession_url, oauth.Appid, oauth.Secretkey, code)
	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if result.Errcode != 0 {
		err = errors.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.Errcode, result.Errmsg)
		return
	}
	return
}
