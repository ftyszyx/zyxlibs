package control

import (
	"encoding/json"

	"github.com/ftyszyx/libs/beego/logs"
	"github.com/ftyszyx/libs/db/mysql"
	"github.com/ftyszyx/libs/models"
	"github.com/pkg/errors"
	"github.com/zyx/shop_server/utils"
)

//通用的查询列表
func (self *BaseController) All() {
	var data = models.AllReqData{And: true}
	err := json.Unmarshal(self.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.Error("%+v", err)
		self.AjaxReturn(utils.ErrorCode, err.Error(), nil)
		return
	}
	self.AllExc(data)
}
func (self *BaseController) AllExc(data models.AllReqData) {
	err, num, datalist := self.AllExcCommon(data, utils.GetAll_type)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	var senddata = make(map[string]interface{})
	senddata["num"] = num
	senddata["list"] = datalist
	self.AjaxReturn(utils.SuccessCode, "", senddata)
}

func (self *BaseController) AllExcCommon(data models.AllReqData, gettype int) (error, int, []mysql.Params) {

	model := models.GetModel(self.control)
	return model.AllExcCommon(self.dboper, model, data, gettype)
}
