package miniprogram

import (
	"github.com/ftyszyx/libs/beego"
	"github.com/ftyszyx/libs/beego/logs"
)

type MiniData struct {
	Appid     string
	Secretkey string
}

func NewInstance(appid string, secretkey string) *MiniData {
	o := new(MiniData)
	o.Appid = appid
	o.Secretkey = secretkey
	return o
}

var Instance *MiniData

func InitMiniProgram() {
	logs.Info("init InitMiniProgram")
	AppID := beego.AppConfig.String("miniprogram.appid")
	AppSecret := beego.AppConfig.String("miniprogram.appsecret")
	Instance = NewInstance(AppID, AppSecret)

}
