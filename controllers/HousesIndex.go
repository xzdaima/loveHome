package controllers

import (
	"github.com/astaxie/beego"
	"loveHome/models"
)

type HousesIndexController struct {
	beego.Controller
}

func (this *HousesIndexController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *HousesIndexController) GetHousesIndex() {
	beego.Info("========== /api/v1.0/houses/index  succ ======")
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)
}
