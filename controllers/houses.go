package controllers

import (
	"encoding/json"
	//  "fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"loveHome/models"

	//"github.com/astaxie/beego/cache"
	//_ "github.com/astaxie/beego/cache/redis"
	"path"
	"strconv"
	"time"
)

type HousesController struct {
	beego.Controller
}

func (this *HousesController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *HousesController) SelectHouse() {
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)

	var aid int
	this.Ctx.Input.Bind(&aid, "aid")
	var sd time.Time
	this.Ctx.Input.Bind(&sd, "sd")
	var ed time.Time
	this.Ctx.Input.Bind(&ed, "ed")
	var p int
	this.Ctx.Input.Bind(&p, "p")
	//now := time.Now()
	if sd.Sub(ed).Seconds() > 0 {
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	if p < 0 {
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}

	/*	cache_conn, err := cache.NewCache("redis", `{"key":"lovehome","conn":"127.0.0.1:6379","dbNum
		":"0"}`)
		if err != nil {
			beego.Info("cache redis conn err,err=", err)
			resp["errno"] = models.RECODE_DBERR
			resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
			return
		}
		areas_info_value := cache_conn.Get("area_info")
		if areas_info_value != nil {
			beego.Info(" ====== get area_info from cache !!! ======")
			var areas_info interface{}
			json.Unmarshal(areas_info_value.([]byte), &areas_info)
			resp["data"] = areas_info
			return
		}*/

	o := orm.NewOrm()
	var house []models.House
	qs := o.QueryTable("house").Filter("area_id", aid)
	num, err := qs.All(&house)
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
	for pos := range house {
		o.LoadRelated(&house[pos], "Area")
		o.LoadRelated(&house[pos], "User")
	}
	var pagecount int
	var onepagehouse []models.House
	var housempcount int
	if int(num)%models.HOUSE_LIST_PAGE_CAPACITY > 0 {
		pagecount = int(num)/models.HOUSE_LIST_PAGE_CAPACITY + 1
		if p < pagecount {
			onepagehouse = house[models.HOUSE_LIST_PAGE_CAPACITY*(p-1) : models.HOUSE_LIST_PAGE_CAPACITY]
			housempcount = models.HOUSE_LIST_PAGE_CAPACITY
		} else {
			onepagehouse = house[models.HOUSE_LIST_PAGE_CAPACITY*(p-1) : int(num)-models.HOUSE_LIST_PAGE_CAPACITY*(p-1)]
			housempcount = int(num) - models.HOUSE_LIST_PAGE_CAPACITY*(p-1)
		}
	} else {
		pagecount = int(num) / models.HOUSE_LIST_PAGE_CAPACITY
		onepagehouse = house[models.HOUSE_LIST_PAGE_CAPACITY*(p-1) : models.HOUSE_LIST_PAGE_CAPACITY]
		housempcount = models.HOUSE_LIST_PAGE_CAPACITY
	}
	housemp := make([]map[string]interface{}, housempcount)

	for pos, value := range onepagehouse {
		housemp[pos] = make(map[string]interface{})
		housemp[pos]["address"] = value.Address
		housemp[pos]["area_name"] = value.Area.Name
		housemp[pos]["ctime"] = value.Ctime
		housemp[pos]["house_id"] = value.Id
		housemp[pos]["img_url"] = "http://192.168.192.134:8080/" + value.Index_image_url
		housemp[pos]["order_count"] = value.Order_count
		housemp[pos]["price"] = value.Price
		housemp[pos]["room_count"] = value.Room_count
		housemp[pos]["title"] = value.Title
		housemp[pos]["user_avatar"] = "http://192.168.192.134:8080/" + value.User.Avatar_url
	}

	datamp := make(map[string]interface{})
	datamp["houses"] = housemp
	datamp["current_page"] = p
	datamp["total_page"] = pagecount
	resp["data"] = datamp

	/*areas_info_str, _ := json.Marshal(areas)
	if err := cache_conn.Put("area_info", areas_info_str, time.Second*3600); err != nil {
		beego.Info("set area_info --> redis fail err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}*/
}

func (this *HousesController) UploadHousePic() {
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
	var id int
	id, _ = strconv.Atoi(this.Ctx.Input.Param(":id"))

	var house models.House
	o := orm.NewOrm()
	res := o.QueryTable("house")
	if err := res.Filter("id", id).One(&house); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	o.LoadRelated(&house, "Images")
	/*if tempuser.Avatar_url != "" {
		models.FDFSDeletFile(tempuser.Avatar_url)
	}*/
	file, header, err := this.GetFile("house_image")
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

	if house.Index_image_url == "" {
		house.Index_image_url = fileId
	}
	houseimage := models.HouseImage{Url: fileId, House: &house}
	house.Images = append(house.Images, &houseimage)

	//	user := models.User{Id: user_id.(int), Avatar_url: fileId}
	if _, err := o.Insert(&houseimage); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	if _, err := o.Update(&house); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	image_url := "http://192.168.192.134:8080/" + fileId
	url_map := make(map[string]interface{})
	url_map["url"] = image_url
	resp["data"] = url_map
	return

}

func (this *HousesController) UploadHouseInfo() {
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
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	if err := qs.Filter("id", userid).One(&user); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	housem := make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &housem)
	beego.Info(housem)
	var area models.Area
	str := housem["area_id"].(string)
	beego.Info(str)
	area.Id, _ = strconv.Atoi(str)

	//beego.`Info
	qs = o.QueryTable("area")
	if err := qs.Filter("id", area.Id).One(&area); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	sinc := housem["facility"].([]interface{})
	beego.Info(sinc)
	facarray := make([]*models.Facility, len(sinc))
	for pos, value := range sinc {
		facarray[pos] = new(models.Facility)
		facarray[pos].Id, _ = strconv.Atoi(value.(string))
	}
	var houseinfo models.House
	//	beego.Info(houseinfo)
	houseinfo.Price, _ = strconv.Atoi(housem["price"].(string))
	houseinfo.Address = housem["address"].(string)
	houseinfo.Room_count, _ = strconv.Atoi(housem["room_count"].(string))
	houseinfo.Unit = housem["unit"].(string)
	houseinfo.Beds = housem["beds"].(string)
	houseinfo.Min_days, _ = strconv.Atoi(housem["min_days"].(string))
	houseinfo.Title = housem["title"].(string)
	houseinfo.Area = &area
	houseinfo.Acreage, _ = strconv.Atoi(housem["acreage"].(string))
	houseinfo.Capacity, _ = strconv.Atoi(housem["capacity"].(string))
	houseinfo.Deposit, _ = strconv.Atoi(housem["deposit"].(string))
	houseinfo.Max_days, _ = strconv.Atoi(housem["max_days"].(string))
	houseinfo.Facilities = facarray
	houseinfo.User = &user

	//m2m:=o.QueryM2M(houseinfo,"facility_houses")
	//m2m.Add()

	house_id, err := o.Insert(&houseinfo)
	if err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	fac := houseinfo.Facilities
	beego.Info(fac)

	for _, value := range fac {
		m2m := o.QueryM2M(value, "Houses")
		_, err := m2m.Add(houseinfo)
		if err != nil {
			resp["errno"] = models.RECODE_DBERR
			resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
			return
		}
	}

	housemap := make(map[string]interface{})
	housemap["house_id"] = house_id
	resp["data"] = housemap
}

func (this *HousesController) GetMyHousesInfo() {
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
	/*	var user models.User
		user.Id = userid.(int)
		o := orm.NewOrm()
		qs := o.QueryTable("user")
		if err := qs.Filter("id", user.Id).One(&user); err != nil {

			resp["errno"] = models.RECODE_DBERR
			resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
			return
		}*/

	o := orm.NewOrm()
	var houses []*models.House
	o.QueryTable("house").Filter("user_id", userid).RelatedSel().All(&houses)
	/*if err != nil {
		fmt.Printf("%d posts read\n", num)
		for _, post := range posts {
			fmt.Printf("Id: %d, UserName: %d, Title: %s\n", post.Id, post.User.UserName, post.Title)
		}
	}*/

	rethouses := make([]map[string]interface{}, len(houses))

	for pos, value := range houses {
		rethouses[pos] = make(map[string]interface{})
		rethouses[pos]["address"] = (*value).Address
		o.QueryTable("area").Filter("id", (*(*value).Area).Id).RelatedSel().One((*value).Area)
		rethouses[pos]["area_name"] = (*(*value).Area).Name
		rethouses[pos]["ctime"] = (*value).Ctime
		rethouses[pos]["house_id"] = (*value).Id
		rethouses[pos]["img_url"] = "http://192.168.192.134:8080/" + (*value).Index_image_url
		rethouses[pos]["order_count"] = (*value).Order_count
		rethouses[pos]["price"] = (*value).Price
		rethouses[pos]["room_count"] = (*value).Room_count
		rethouses[pos]["title"] = (*value).Title
		o.QueryTable("user").Filter("id", (*value).User.Id).RelatedSel().One((*value).User)
		rethouses[pos]["user_avatar"] = "http://192.168.192.134:8080/" + (*value).User.Avatar_url
	}

	Housemap := make(map[string]interface{})
	Housemap["houses"] = rethouses
	resp["data"] = Housemap
}

func (this *HousesController) GetHouse(hid int, resp map[string]interface{}) map[string]interface{} {

	var house models.House
	o := orm.NewOrm()
	if err := o.QueryTable("house").Filter("id", hid).RelatedSel().One(&house); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return nil
	}
	o.LoadRelated(&house, "Facilities")
	o.LoadRelated(&house, "Images")
	//	house.Facilities
	beego.Info(house)
	/*	beego.Info(house)
		beego.Info(*(house.Area))
		beego.Info(*(house.User))
		beego.Info(*(house.Images))*/
	//com := make([]interface{}, 1)
	var com []interface{}

	fac := make([]int, len(house.Facilities))
	for pos, value := range house.Facilities {
		fac[pos] = (*value).Id
	}
	imaul := make([]string, len(house.Images))
	for pos, value := range house.Images {
		imaul[pos] = "http://192.168.192.134:8080/" + (*value).Url
	}
	var user models.User
	if err := o.QueryTable("user").Filter("id", (*house.User).Id).One(&user); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return nil
	}

	housedatamap := make(map[string]interface{})
	housedatamap["acreage"] = house.Acreage
	housedatamap["address"] = house.Address

	housedatamap["beds"] = house.Beds
	housedatamap["capacity"] = house.Capacity

	housedatamap["comments"] = com
	housedatamap["deposit"] = house.Deposit
	housedatamap["facilities"] = fac
	housedatamap["hid"] = house.Id
	housedatamap["img_urls"] = imaul
	housedatamap["max_days"] = house.Max_days
	housedatamap["min_days"] = house.Min_days
	housedatamap["price"] = house.Price
	housedatamap["room_count"] = house.Room_count
	housedatamap["title"] = house.Title
	housedatamap["unit"] = house.Unit
	housedatamap["user_avatar"] = "http://192.168.192.134:8080/" + user.Avatar_url
	housedatamap["user_id"] = user.Id
	housedatamap["user_name"] = user.Name

	return housedatamap
}

func (this *HousesController) GetHouseInfo() {
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)
	var id int
	id, _ = strconv.Atoi(this.Ctx.Input.Param(":id"))
	userid := this.GetSession("user_id")
	housedatamap := this.GetHouse(id, resp)
	//	housedatamap := make(map[string]interface{})
	//	housedatamap["acreage"]
	//	housedatamap["address"]
	hou := make(map[string]interface{})
	//	housedatamap["beds"]
	//	housedatamap["capacity"]
	hou["house"] = housedatamap
	hou["user_id"] = userid
	resp["data"] = hou
}
