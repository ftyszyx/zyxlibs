package kuaidiniao

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
	"github.com/ftyszyx/libs/beego/httplib"
	"github.com/ftyszyx/libs/beego/logs"
	zyxstr "github.com/ftyszyx/libs/string"
	"github.com/pkg/errors"
)

type KuaiData struct {
	AcceptTime    string
	AcceptStation string
}

type KuaiResp struct {
	EBusinessID  string
	OrderCode    string
	ShipperCode  string
	LogisticCode string
	Success      bool
	Reason       string
	State        string
	Traces       []KuaiData
}

type SendParam struct {
	OrderCode    string `json:"OrderCode"`
	ShipperCode  string `json:"ShipperCode"`
	LogisticCode string `json:"LogisticCode"`
}

var LogisticsCodeArr = map[string]string{
	"huitongkuaidi": "HTKY",
	"ems":           "EMS",
	"shunfeng":      "SF",
	"tiantian":      "HHTT",
	"yuantong":      "YTO",
	"yunda":         "YD",
	"zhongtong":     "ZTO"
}


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

//获取物流信息
func GetKuaiInfo(costomerid string, key string, company string, num string) (*KuaiResp, error) {
	var param SendParam
	param.OrderCode = ""
	param.ShipperCode = LogisticsCodeArr[company]
	param.LogisticCode = num
	parambuf, err := json.Marshal(param)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var sigContent = string(parambuf) + key
	logs.Info("sigContent:%s", sigContent)
	signstr := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(zyxstr.GetStrMD5(sigContent))))

	urlstr := "http://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
	req := httplib.Post(urlstr)
	sendata := make(map[string]interface{})
	sendata["EBusinessID"] = costomerid
	sendata["RequestType"] = "1002"
	sendata["RequestData"] = url.QueryEscape(string(parambuf))
	sendata["DataType"] = "2"
	sendata["DataSign"] = signstr

	logs.Info("send data:%v", sendata)

	for key, value := range sendata {
		req.Param(key, fmt.Sprintf("%+v", value))
	}
	req.Header("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	var respdata []byte
	respdata, err = req.Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	getData := new(KuaiResp)
	//logs.Info("get data:%s", string(respdata))
	err = json.Unmarshal(respdata, getData)
	if err != nil {
		logs.Info("parse data err")
		return nil, errors.WithStack(err)
	}
	if getData.Success == false {
		return nil, errors.New(getData.Reason)
	}
	return getData, nil
}
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
 
const g_encrypt_type="url"
const g_version="1.0"
const g_format="json"

const ERR_SERVER_ERR=505 //系统异常 
const SUCCESS=200//成功
const ERR_BADREQUEST=400 //错误请求
const ERR_BADMETHOD=405//禁用的方法(不支持服务名)
const ERR_SIGNERROR=420//签名验证失败
const ERR_BADREQUEST2=420//请求格式错误【参数名】


func Getparam(costomerid string, method string,curtime string,sign string) string {
	//curtime := time.Now().Unix()
	url := fmt.Sprintf("?partnerid=%s&timestamp=%s&method=%s&sign=%s&format=%s&encrypt=%s&version=%s",
	costomerid, curtime,method,sign,g_format,g_encrypt_type,g_version)
	return url
}

//获取签名
func GetSign(data string,method string, costomerid string,curtime string,key string ) string {
	datastr := fmt.Sprintf("data=%sencrypt=%format=%smethod=%partnerid=%timestamp=%sversion=%s%s",
	 data, g_encrypt_type, g_format, method, costomerid, curtime,g_version,key)
	 signsstr := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(zyxstr.GetStrMD5(datastr))))
	 return signsstr
}