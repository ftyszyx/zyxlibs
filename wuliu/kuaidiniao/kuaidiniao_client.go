package kuaidiniao

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

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

type SendQueryParam struct {
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
	"zhongtong":     "ZTO"}

var LogisticsCodeArr2 = map[string]string{
	"HTKY": "huitongkuaidi",
	"EMS":  "ems",
	"SF":   "shunfeng",
	"HHTT": "tiantian",
	"YTO":  "yuantong",
	"YD":   "yunda",
	"ZTO":  "zhongtong"}

//push相关的结构
type PushTrace struct {
	AcceptTime    string
	AcceptStation string
	Remark        string
}

type PushDataInfo struct {
	EBusinessID           string
	ShipperCode           string
	LogisticCode          string
	Success               bool
	Reason                string
	EstimatedDeliveryTime string
	Traces                []PushTrace
}

type PushRequestData struct {
	EBusinessID string
	PushTime    string
	Count       string
	Data        []PushDataInfo
}

type PushRespData struct {
	EBusinessID string
	UpdateTime  string
	Success     bool
}

//获取要发送的结构
func Client_GetSendData(paramstr string, key string, costomerid string, cmd string) map[string]interface{} {

	signstr := Client_GetSign(paramstr, key)
	sendata := make(map[string]interface{})
	sendata["EBusinessID"] = costomerid
	sendata["RequestType"] = cmd
	sendata["RequestData"] = url.QueryEscape(paramstr)
	sendata["DataType"] = "2"
	sendata["DataSign"] = signstr
	logs.Info("send data:%v", sendata)
	return sendata
}

func Client_GetSign(paramstr string, key string) string {
	var sigContent = paramstr + key
	logs.Info("sigContent:%s", sigContent)
	//signstr := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(zyxstr.GetStrMD5(sigContent))))
	signstr := base64.StdEncoding.EncodeToString([]byte(zyxstr.GetStrMD5(sigContent)))
	return signstr
}

//查询物流信息
func Client_Query(costomerid string, key string, sendparam SendQueryParam) (*KuaiResp, error) {
	parambuf, err := json.Marshal(sendparam)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	urlstr := "http://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
	req := httplib.Post(urlstr)
	sendata := Client_GetSendData(string(parambuf), key, costomerid, "1002")
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

type Receiver struct {
	Name         string `json:"Name"`
	Tel          string `json:"Tel"`
	Mobile       string `json:"Mobile"`
	ProvinceName string `json:"ProvinceName"`
	CityName     string `json:"CityName"`
	ExpAreaName  string `json:"ExpAreaName"`
	Address      string `json:"Address"`
}

type Sender struct {
	Name         string `json:"Name"`
	Tel          string `json:"Tel"`
	Mobile       string `json:"Mobile"`
	ProvinceName string `json:"ProvinceName"`
	CityName     string `json:"CityName"`
	ExpAreaName  string `json:"ExpAreaName"`
	Address      string `json:"Address"`
}

type Addlister_SendParam struct {
	ShipperCode   string   `json:"ShipperCode"`
	LogisticCode  string   `json:"LogisticCode"`
	Sender_info   Sender   `json:"Sender"`
	Receiver_info Receiver `json:"Receiver"`
}

type Listener_resp struct {
	EBusinessID string
	UpdateTime  string
	Success     bool
	Reason      string
}

//订阅
func Client_Addlistner(costomerid string, key string, sendparam Addlister_SendParam) error {

	parambuf, err := json.Marshal(sendparam)
	if err != nil {
		return errors.WithStack(err)
	}
	urlstr := "https://api.kdniao.com/api/dist"
	//urlstr := "http://sandboxapi.kdniao.com:8080/kdniaosandbox/gateway/exterfaceInvoke.json"
	req := httplib.Post(urlstr)
	sendata := Client_GetSendData(string(parambuf), key, costomerid, "1008")
	logs.Info("sendata:%v", sendata)
	for key, value := range sendata {
		req.Param(key, fmt.Sprintf("%+v", value))
	}
	req.Header("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	var respdata []byte
	respdata, err = req.Bytes()
	if err != nil {
		return errors.WithStack(err)
	}
	getData := new(Listener_resp)
	logs.Info("get data:%s", string(respdata))
	err = json.Unmarshal(respdata, getData)
	if err != nil {
		logs.Info("parse data err")
		return errors.WithStack(err)
	}
	if getData.Success == false {
		return errors.New(getData.Reason)
	}
	return nil
}
