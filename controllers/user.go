package controllers

import (
	"encoding/json"
	//	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"loveHome/models"
	"path"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *UserController) UpdataUserName() {
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)
	userid := this.GetSession("user_id")
	if userid == nil {
		resp["errno"] = models.RECODE_LOGINERR
		resp["errmsg"] = models.RecodeText(models.RECODE_LOGINERR)
		return
	}
	var usernamemap map[string]interface{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &usernamemap)
	//	beego.Info(usernamemap)
	var user models.User
	user.Id = userid.(int)
	user.Name = usernamemap["name"].(string)
	o := orm.NewOrm()
	if _, err := o.Update(&user, "name"); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	resp["data"] = usernamemap

}

func (this *UserController) GetUserInfo() {

	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)
	userid := this.GetSession("user_id")
	if userid == nil {
		resp["errno"] = models.RECODE_LOGINERR
		resp["errmsg"] = models.RecodeText(models.RECODE_LOGINERR)
		return
	}

	var user models.User
	user.Id = userid.(int)

	o := orm.NewOrm()
	qs := o.QueryTable("user")
	if err := qs.Filter("id", user.Id).One(&user); err != nil {

		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	resp["data"] = user

}
func (this *UserController) GetUserRealInfo() {

	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)
	userid := this.GetSession("user_id")
	if userid == nil {
		resp["errno"] = models.RECODE_LOGINERR
		resp["errmsg"] = models.RecodeText(models.RECODE_LOGINERR)
		return
	}

	var user models.User
	user.Id = userid.(int)

	o := orm.NewOrm()
	qs := o.QueryTable("user")
	if err := qs.Filter("id", user.Id).One(&user); err != nil {

		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	realusermap := make(map[string]interface{})
	realusermap["user_id"] = user.Id
	realusermap["name"] = user.Name
	realusermap["password"] = user.Password_hash
	realusermap["mobile"] = user.Mobile
	realusermap["real_name"] = user.Real_name
	realusermap["id_card"] = user.Id_card
	realusermap["avatar_url"] = user.Avatar_url

	resp["data"] = realusermap

}
func (this *UserController) UpdataUserRealInfo() {
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)
	userid := this.GetSession("user_id")
	if userid == nil {
		resp["errno"] = models.RECODE_LOGINERR
		resp["errmsg"] = models.RecodeText(models.RECODE_LOGINERR)
		return
	}
	var usernamemap map[string]interface{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &usernamemap)
	//	beego.Info(usernamemap)
	var user models.User
	user.Id = userid.(int)
	user.Real_name = usernamemap["real_name"].(string)
	user.Id_card = usernamemap["id_card"].(string)
	o := orm.NewOrm()
	if _, err := o.Update(&user, "id_card", "real_name"); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//	resp["data"] = usernamemap

}
func (this *UserController) Reg() {
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	var regRequstMap = make(map[string]interface{})

	json.Unmarshal(this.Ctx.Input.RequestBody, &regRequstMap)
	//	fmt.Println(regRequstMap)
	beego.Info("mobile=", regRequstMap["mobile"])
	beego.Info("password=", regRequstMap["password"])
	beego.Info("sms_code=", regRequstMap["sms_code"])

	if regRequstMap["mobile"] == "" || regRequstMap["password"] == "" || regRequstMap["sms_code"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}
	user := models.User{}
	user.Mobile = regRequstMap["mobile"].(string)
	user.Password_hash = regRequstMap["password"].(string)
	user.Name = regRequstMap["mobile"].(string)

	o := orm.NewOrm()

	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("insert error=", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	beego.Info("reg succ!!!user id=", id)
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", id)
	this.SetSession("mobile", user.Mobile)

	return

}

func (this *UserController) Login() {
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)

	var loginRequestMap = make(map[string]interface{})

	json.Unmarshal(this.Ctx.Input.RequestBody, &loginRequestMap)
	beego.Info("mobile = ", loginRequestMap["mobile"])
	beego.Info("password = ", loginRequestMap["password"])

	if loginRequestMap["mobile"] == "" || loginRequestMap["password"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}
	var user models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	if err := qs.Filter("mobile", loginRequestMap["mobile"]).One(&user); err != nil {
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}
	if user.Password_hash != loginRequestMap["password"].(string) {
		resp["errno"] = models.RECODE_PWDERR
		resp["errmsg"] = models.RecodeText(models.RECODE_PWDERR)
		return
	}

	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)
}

func (this *UserController) UploadAvatar() {
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)

	user_id := this.GetSession("user_id")
	if user_id == nil {
		resp["errno"] = models.RECODE_LOGINERR
		resp["errmsg"] = models.RecodeText(models.RECODE_LOGINERR)
		return
	}

	file, header, err := this.GetFile("avatar")
	if err != nil {
		resp["errno"] = models.RECODE_SERVERERR

		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	fileBuffer := make([]byte, header.Size)
	if _, err := file.Read(fileBuffer); err != nil {
		resp["errno"] = models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		return
	}
	suffix := path.Ext(header.Filename)
	groupName, fileId, err := models.FDFSUploadByBuffer(fileBuffer, suffix)
	if err != nil {
		resp["errno"] = models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		beego.Info("upload file to fastdfs error err = ", err)
		return
	}
	beego.Info("fdfs upload succ groupname=", groupName, "fileid=", fileId)

	user := models.User{Id: user_id.(int), Avatar_url: fileId}

	var tempuser models.User
	o := orm.NewOrm()
	res := o.QueryTable("user")
	res.Filter("id", user.Id).One(&tempuser)
	if tempuser.Avatar_url != "" {
		models.FDFSDeletFile(tempuser.Avatar_url)
	}

	if _, err := o.Update(&user, "avatar_url"); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	avatar_url := "http://192.168.192.134:8080/" + fileId
	url_map := make(map[string]interface{})
	url_map["avatar_url"] = avatar_url
	resp["data"] = url_map
	return
}
