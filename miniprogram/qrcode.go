package miniprogram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ftyszyx/libs/beego"
	"github.com/ftyszyx/libs/beego/httplib"
	"github.com/ftyszyx/libs/qiniu"
	zyxstr "github.com/ftyszyx/libs/string"
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

func (oauth *MiniData) GetQrCode(code string, scene string, page string, width string) (url string, err error) {

	keystr := page + scene + width
	md5str := zyxstr.GetStrMD5(keystr)
	tempfolder := beego.AppConfig.String("site.tempfolder")
	host := beego.AppConfig.String("qiniu.host")
	fileName := md5str + ".jpg"
	filepath := tempfolder + fileName
	url = host + fileName

	//判断是否有同样的文件
	bucketname := beego.AppConfig.String("qiniu.bucket")
	manger := qiniu.GetManger()
	_, fileerr := manger.Stat(bucketname, url)

	if fileerr == nil {
		//说明文件存在
		return
	}

	//获取token
	res, tokenerr := oauth.GetToken(code)
	if tokenerr != nil {
		err = tokenerr
	}

	//请求码
	urlstr := fmt.Sprintf(getqrcode_url, res.Access_token)
	var response []byte
	req := httplib.Post(urlstr)
	//req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	//req.Header("Content-Type", "application/x-www-form-urlencoded")

	req.Param("scene", scene)
	req.Param("page", page)
	req.Param("width", width)
	// req.Param("access_token", access_token)
	response, err = req.Bytes()
	if err != nil {
		return
	}
	var getData map[string]interface{}
	err = json.Unmarshal(response, &getData)
	if err != nil {
		return
	}
	if getData["errcode"] != nil {
		err = errors.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", getData["errcode"], getData["errmsg"])
		return
	}

	//写入临时文件
	err = ioutil.WriteFile(filepath, response, 0755)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	//上传七牛
	_, err = qiniu.UploadFile(fileName, filepath, bucketname)

	return
}
