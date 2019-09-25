package kuaidiniao

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ftyszyx/libs/beego/httplib"
	"github.com/ftyszyx/libs/beego/logs"
	zyxstr "github.com/ftyszyx/libs/string"
	"github.com/pkg/errors"
)

//给快递鸟提供服务

/*
querytrace 路由信息查询
registerorder 订单注册
pushtrace 路由信息推送
pushorderstate 订单状态推送
agentpointtrace 末端路由下发
agentpushtrace 末端路由推送
querysite 网点查询
applyservice 面单账号申请
serviceresult 面单账号申请反馈
electronicorder 面单下单
querybalance 余量查询
recyclecode 面单号回收
createorder 预约取件
Cancleorder 订单取消
createrealname 个人实名信息上传
updaterealname 更新个人实名信息
canclerealname 删除个人实名信息
*/

/*
partnerid string R 合作商 ID
timestamp string R 时间戳，格式：1498489639
method string R 服务名，参考服务名定义
sign string R 数据签名，参考数据签名
format string R 报文格式：json,xml；默认 json
encrypt string R 报文(data)加密方式：url(utf-8),aes；
默认 url(utf-8)
version string R 版本号，默认 1.0
*/

const g_encrypt_type = "url"
const g_version = "1.0"
const g_format = "json"

const ERR_SERVER_ERR = 505  //系统异常
const SUCCESS = 200         //成功
const ERR_BADREQUEST = 400  //错误请求
const ERR_BADMETHOD = 405   //禁用的方法(不支持服务名)
const ERR_SIGNERROR = 420   //签名验证失败
const ERR_BADREQUEST2 = 420 //请求格式错误【参数名】

func SrvGetparam(costomerid string, method string, curtime string, sign string) string {
	//curtime := time.Now().Unix()
	url := fmt.Sprintf("?partnerid=%s&timestamp=%s&method=%s&sign=%s&format=%s&encrypt=%s&version=%s",
		costomerid, curtime, method, sign, g_format, g_encrypt_type, g_version)
	return url
}

//获取签名
func SrvGetSign(data string, method string, costomerid string, curtime string, key string) string {
	datastr := fmt.Sprintf("data=%sencrypt=%sformat=%smethod=%spartnerid=%stimestamp=%sversion=%s%s",
		data, g_encrypt_type, g_format, method, costomerid, curtime, g_version, key)
	logs.Info("datastr:%s", datastr)

	//signsstr := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(zyxstr.GetStrMD5(datastr))))
	signsstr := strings.ToLower(base64.StdEncoding.EncodeToString([]byte(zyxstr.GetStrMD5(datastr))))
	return signsstr
}

type Srv_PushData struct {
	OrderCode   string
	WaybillCode string
	ScanType    string //PRIORA
	ScanData    string
	TraceDesc   string
	CallBack    string
}

type Srv_PushDataResp struct {
	PartnerId  string
	Success    string
	ResultCode string
	Reason     string
}

//发送推送
func SrvPush(tracelist []Srv_PushData, costomerid string, key string) error {
	jsonstr, err := json.Marshal(tracelist)
	if err != nil {
		return errors.WithStack(err)
	}
	curtime := time.Now().Unix()
	timestr := strconv.FormatInt(curtime, 10)
	signstr := SrvGetSign(string(jsonstr), "pushtrace", costomerid, timestr, key)
	paramstr := SrvGetparam(costomerid, "pushtrace", timestr, signstr)
	urlstr := "http://183.62.170.46:38093" + paramstr
	req := httplib.Post(urlstr)
	req.Body(string(jsonstr))
	req.Header("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	var respdata []byte
	respdata, err = req.Bytes()
	if err != nil {
		return errors.WithStack(err)
	}
	getData := new(Srv_PushDataResp)
	//logs.Info("get data:%s", string(respdata))
	err = json.Unmarshal(respdata, getData)
	if err != nil {
		return errors.WithStack(err)
	}
	if getData.Success == "false" {
		return errors.New(getData.Reason)
	}
	return nil

}
