package controllers

import (
	//	"fmt"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"loveHome/models"
	"time"
)

type AreasController struct {
	beego.Controller
}

/*func (c *AreasController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}*/
func (this *AreasController) RetData(resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *AreasController) GetAreaInfo() {
	beego.Info("--------------GetContro-------")
	resp := make(map[string]interface{})
	resp["errno"] = 0
	resp["errmsg"] = "OK"
	defer this.RetData(resp)

	cache_conn, err := cache.NewCache("redis", `{"key":"lovehome","conn":"127.0.0.1:6379","dbNum":"0"}`)
	if err != nil {
		beego.Info("cache redis conn err,err=", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//	cache_conn.Put("cccc", "----lala---", time.Second*300)

	//	value := cache_conn.Get("cccc")
	//fmt.Printf("%s", value)
	//	if value != nil {
	//		fmt.Printf("%s", value)
	//	}

	areas_info_value := cache_conn.Get("area_info")
	if areas_info_value != nil {
		beego.Info(" ====== get area_info from cache !!! ======")
		var areas_info interface{}
		json.Unmarshal(areas_info_value.([]byte), &areas_info)
		resp["data"] = areas_info
		return
	}

	o := orm.NewOrm()
	var areas []models.Area
	qs := o.QueryTable("area")
	num, err := qs.All(&areas)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = "查询数据库失败"
		return
	}
	if num == 0 {
		resp["errno"] = 4002
		resp["errmsg"] = "没有数据"
		return
	}
	resp["data"] = areas

	areas_info_str, _ := json.Marshal(areas)
	if err := cache_conn.Put("area_info", areas_info_str, time.Second*3600); err != nil {
		beego.Info("set area_info --> redis fail err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	return

}
