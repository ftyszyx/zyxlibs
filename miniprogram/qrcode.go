package miniprogram

import (
	"encoding/json"
	"fmt"

	"github.com/ftyszyx/libs/beego/httplib"
	"github.com/pkg/errors"
)

const (
	getqrcode_url = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)

type QRCode struct {
	ContentType string `json:"contentType"`
	Buffer      []byte `json:"buffer"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

func (oauth *MiniData) GetQrCode(access_token string, scene string, page string, width string) (result QRCode, err error) {
	urlstr := fmt.Sprintf(getqrcode_url, access_token)
	var response []byte
	req := httplib.Post(urlstr)
	req.Param("scene", scene)
	req.Param("page", page)
	req.Param("width", width)
	// req.Param("access_token", access_token)
	response, err = req.Bytes()
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
