package alibaba

import (
"fmt"
	"github.com/pkg/errors"
	"github.com/ftyszyx/libs/beego/logs"
"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

var Templatecode string
var SignName string
var APPID string
var APPKEY string

func SendCode(phone string, code string) error{
	logs.Info("phoneis %s code:%s\n",phone, code)
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", APPID, APPKEY)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phone
	request.SignName = SignName
	request.TemplateCode = Templatecode
	request.TemplateParam = fmt.Sprintf("{\"code\":\"%s\"}",code)

	response, err := client.SendSms(request)
	if err != nil {
		return  errors.New(err.Error())
	}
	logs.Info("response is %#v\n", response)
	if response.Code!="OK"{
		return  errors.New(response.Message)
	}
	return nil
}

