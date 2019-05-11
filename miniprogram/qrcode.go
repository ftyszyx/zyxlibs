package miniprogram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ftyszyx/libs/beego"
	"github.com/ftyszyx/libs/beego/logs"
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

type QRCoderReq struct {
	Page string `json:"page,omitempty"`
	// path 识别二维码后进入小程序的页面链接
	Path string `json:"path,omitempty"`
	// width 图片宽度
	Width int `json:"width,omitempty"`
	// scene 参数数据
	Scene string `json:"scene,omitempty"`
	// autoColor 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调
	AutoColor bool `json:"auto_color,omitempty"`
	// lineColor AutoColor 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"},十进制表示
	LineColor Color `json:"line_color,omitempty"`
	// isHyaline 是否需要透明底色
	IsHyaline bool `json:"is_hyaline,omitempty"`
}

type Color struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

func (oauth *MiniData) GetQrCode(scene string, page string, width int) (url string, err error) {

	keystr := page + scene + strconv.Itoa(width)
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
	res, tokenerr := oauth.GetToken()
	if tokenerr != nil {
		err = tokenerr
	}

	//请求码
	urlstr := fmt.Sprintf(getqrcode_url, res.Access_token)

	senddata := new(QRCoderReq)
	senddata.Page = page
	senddata.Scene = scene
	senddata.Width = width
	reqbody, jsonerr := json.Marshal(senddata)
	if jsonerr != nil {
		err = jsonerr
		return
	}
	qrres, qrerr := http.Post(urlstr, "application/json", strings.NewReader(string(reqbody)))
	if qrerr != nil {
		err = qrerr
		return
	}

	logs.Info("url:%s ", urlstr)

	defer qrres.Body.Close()
	var response []byte
	response, err = ioutil.ReadAll(qrres.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	//logs.Info("response:%s ", string(response))
	var getData map[string]interface{}

	err = json.Unmarshal(response, &getData)
	if err != nil {
		err = errors.WithStack(err)
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
