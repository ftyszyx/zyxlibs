package kuaidiniao

import (
	"encoding/base64"
	"encoding/json"
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

type SendParam struct {
	OrderCode    string `json:"OrderCode"`
	ShipperCode  string `json:"ShipperCode"`
	LogisticCode string `json:"LogisticCode"`
}

func GetKuaiInfo(costomerid string, key string, company string, num string) (*KuaiResp, error) {
	var param SendParam
	param.OrderCode = ""
	param.ShipperCode = company
	param.LogisticCode = num
	parambuf, err := json.Marshal(param)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var sigContent = string(parambuf) + key
	signstr := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(zyxstr.GetStrMD5(sigContent))))

	urlstr := "http://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
	req := httplib.Post(urlstr)
	sendata := make(map[string]interface{})
	sendata["EBusinessID"] = costomerid
	sendata["RequestType"] = "1002"
	sendata["RequestData"] = url.QueryEscape(string(parambuf))
	sendata["DataType"] = "2"
	sendata["DataSign"] = signstr

	reqbuf, err := json.Marshal(sendata)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	logs.Info("send data:%s", string(reqbuf))

	req.Body(string(reqbuf))
	req.Header("Content-Type", "application/x-www-form-urlencoded")

	var respdata []byte
	respdata, err = req.Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	getData := new(KuaiResp)
	logs.Info("get data:%s", string(respdata))
	err = json.Unmarshal(respdata, getData)
	if err != nil {
		logs.Info("parse data err")
		return nil, errors.WithStack(err)
	}
	// logs.Info("get data:%v", getData)
	if getData.Success == false {
		return nil, errors.New(getData.Reason)
	}
	return getData, nil
}
