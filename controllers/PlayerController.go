package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	//"github.com/chnzrb/admin/proto"
	//"github.com/chnzrb/myadmin/proto"
	//"github.com/chnzrb/myadmin/proto"
	"github.com/chnzrb/myadmin/proto"
	"time"
)

type PlayerController struct {
	BaseController
}

func (c *PlayerController) List() {
	var params models.PlayerQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询用户列表:%+v", params)
	data, total := models.GetPlayerList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取玩家列表成功", result)
}

func (c *PlayerController) Detail() {
	var params struct {
		PlatformId int    `json:"platformId"`
		ServerId   string `json:"serverId"`
		PlayerId   int    `json:"playerId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询玩家详细信息:%+v", params)
	playerDetail, err := models.GetPlayerDetail(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err, "查询玩家详细信息失败")
	c.Result(enums.CodeSuccess, "获取玩家详细信息成功", playerDetail)
}

func (c *PlayerController) One() {
	//var params struct {
	//	PlatformId int `json:"platformId"`
	//	ServerId string `json:"serverId"`
	//	PlayerId int `json:"playerId"`
	//}
	//err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//utils.CheckError(err)
	//logs.Info("查询玩家详细信息:%+v", params)
	platformId, err := c.GetInt("platformId")
	c.CheckError(err)
	serverId := c.GetString("serverId")
	playerName := c.GetString("playerName")
	c.CheckError(err)
	player, err := models.GetPlayerByPlatformIdAndSidAndNickname(platformId, serverId, playerName)
	c.CheckError(err, "查询玩家失败")
	c.Result(enums.CodeSuccess, "获取玩家成功", player)
}

func (c *PlayerController) PlayerLoinLogList() {
	var params models.PlayerLoginLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询玩家详细信息:%+v", params)
	data, total := models.GetPlayerLoginLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取玩家登录日志", result)
}

func (c *PlayerController) PlayerOnlineLogList() {
	var params models.PlayerOnlineLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询在线日志:%+v", params)
	data, total := models.GetPlayerOnlineLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取在线日志", result)
}

func (c *PlayerController) MailLogList() {
	var params models.MailLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询邮件日志:%+v", params)
	data, total := models.GetMailLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取邮件日志", result)
	v := [] string {
		"d",
		"1",
		"2",
		"d",
	}
	utils.RemoveDuplicateArray(v)
}

func (c *PlayerController) DelMailLog() {
	var idList []int
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &idList)
	utils.CheckError(err)
	logs.Info("删除邮件列表:%+v", idList)
	err = models.DeleteMailLog(idList)
	c.CheckError(err, "删除邮件失败")
	c.Result(enums.CodeSuccess, "成功删除邮件", idList)
}
func (c *PlayerController) GetServerGeneralize() {
	var params models.ServerGeneralizeQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询服务器概况:%+v", params)
	data, err := models.GetServerGeneralize(params.PlatformId, params.ServerId)
	c.CheckError(err)
	c.Result(enums.CodeSuccess, "获取服务器概况", data)
}

func (c *PlayerController) ForbidLogList() {
	var params models.ForbidLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询封禁日志:%+v", params)
	data, total := models.GetForbidLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取封禁日志", result)
}

func (c *PlayerController) SetDisable() {
	var params struct {
		PlatformId int
		ServerId   string
		PlayerName string
		Type       int32
		Sec        int32
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("封禁:%+v", params)
	//platformId, err := c.GetInt("platformId")
	//c.CheckError(err)
	//serverId:= c.GetString("serverId")
	//playerId, err := c.GetInt("playerId")
	//c.CheckError(err)
	conn, err := models.GetWsByPlatformIdAndSid(params.PlatformId, params.ServerId)
	c.CheckError(err)
	defer conn.Close()
	player, err := models.GetPlayerByPlatformIdAndSidAndNickname(params.PlatformId, params.ServerId, params.PlayerName)
	c.CheckError(err)
	request := gm.MSetDisableTos{Token: proto.String(""), Type: proto.Int32(params.Type), PlayerId: proto.Int32(int32(player.Id)), Sec: proto.Int32(params.Sec)}
	mRequest, err := proto.Marshal(&request)
	c.CheckError(err)

	_, err = conn.Write(Packet(9901, mRequest))
	c.CheckError(err)
	var receive = make([]byte, 100, 100)
	n, err := conn.Read(receive)
	c.CheckError(err)
	respone := &gm.MSetDisableToc{}
	//f :=receive[1]
	//isZip := f >> 7 == 1
	//fmt.Printf("isZip:%v", isZip)
	data := receive[5:n]
	//if isZip{
	//	data = DoZlibUnCompress(data)
	//}
	err = proto.Unmarshal(data, respone)
	c.CheckError(err)

	if *respone.Result == gm.MSetDisableToc_success {
		var forbidTime int32
		if params.Sec > 0 {
			forbidTime = int32(time.Now().Unix()) + params.Sec
		} else {
			forbidTime = 0
		}
		mailLog := &models.ForbidLog{
			PlatformId: params.PlatformId,
			ServerId:   string(params.ServerId),
			PlayerName: params.PlayerName,
			ForbidType: params.Type,
			ForbidTime: forbidTime,
			Time:       time.Now().Unix(),
			UserId:     c.curUser.Id,
		}
		err = models.Db.Save(&mailLog).Error
		c.CheckError(err, "写封禁日志失败")
		c.Result(enums.CodeSuccess, "封禁成功", 0)
	} else {
		c.Result(enums.CodeFail, "封禁失败", 0)
	}
	//conn.Read()

}

func (c *PlayerController) SendMail() {
	type platformServer struct {
		PlatformId   int
		ServerIdList [] string
	}
	var params struct {
		PlatformId         int
		ServerIdList       [] string
		platformServerList [] *platformServer
		PlayerNameList     string
		MailItemList       [] *gm.MSendMailTosProp
		Title              string
		Content            string
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("发送邮件:%+v", params)
	logs.Info("发送邮件:%+v", params.MailItemList)
	//platformId, err := c.GetInt("platformId")
	//c.CheckError(err)
	//serverId:= c.GetString("serverId")
	//playerId, err := c.GetInt("playerId")
	//c.CheckError
	serverIdList, err := json.Marshal(params.ServerIdList)
	c.CheckError(err)
	itemList, err := json.Marshal(params.MailItemList)
	c.CheckError(err)
	mailLog := &models.MailLog{
		PlatformId:     params.PlatformId,
		ServerIdList:   string(serverIdList),
		Title:          params.Title,
		Content:        params.Content,
		Time:           time.Now().Unix(),
		UserId:         c.curUser.Id,
		ItemList:       string(itemList),
		PlayerNameList: params.PlayerNameList,
		Status:         0,
	}
	err = models.Db.Save(&mailLog).Error
	c.CheckError(err, "写邮件日志失败")
	for _, serverId := range params.ServerIdList {
		conn, err := models.GetWsByPlatformIdAndSid(params.PlatformId, serverId)
		c.CheckError(err)
		defer conn.Close()
		request := gm.MSendMailTos{
			Token:          proto.String(""),
			Title:          proto.String(params.Title),
			Content:        proto.String(params.Content),
			PlayerNameList: proto.String(params.PlayerNameList),
			PropList:       params.MailItemList,
		}
		mRequest, err := proto.Marshal(&request)
		c.CheckError(err)

		_, err = conn.Write(Packet(9903, mRequest))
		c.CheckError(err)
		var receive = make([]byte, 100, 100)
		n, err := conn.Read(receive)
		c.CheckError(err)
		respone := &gm.MSendMailToc{}
		data := receive[5:n]
		err = proto.Unmarshal(data, respone)
		c.CheckError(err)

		if *respone.Result == gm.MSendMailToc_success {
			logs.Info("发送邮件成功:+v", request)
		} else {
			c.Result(enums.CodeFail, "发送邮件失败", 0)
		}
	}
	c.Result(enums.CodeSuccess, "发送邮件成功", 0)
	//conn.Read()
}

//封包
func Packet(methodNum int, message []byte) []byte {
	return append(append([]byte{0}, IntToBytes(methodNum)...), message...)
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
