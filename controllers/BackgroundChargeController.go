// 后台充值
package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/chnzrb/myadmin/enums"
	"net/http"
	"io/ioutil"
	"fmt"
)

type BackgroundController struct {
	BaseController
}

func (c *BackgroundController) List() {
	var params models.BackgroundChargeLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询后台充值日志列表:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetBackgroundChargeLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取后台充值日志列表成功", result)
}

func (c *BackgroundController) Charge() {
	var params struct {
		Account     string
		Ip          string
		PlayerId    int
		PlatformId  string
		ChargeValue int
		ServerId    string
		ChargeType  string
		ItemId      int
	}
	var result struct {
		Code    int
		Message string
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("后台充值:%v", params)
	c.CheckError(err)

	player, err := models.GetPlayerOne(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err)

	accountType := models.GetAccountType(params.PlatformId, player.AccId)

	if accountType != 1 {
		c.Result(enums.CodeFail, fmt.Sprintf("后台充值失败: %s.%s 不是内部帐号", params.ServerId, player.Nickname), 0)
	}

	args := fmt.Sprintf("player_id=%d&game_charge_id=0&charge_item_id=%d&item_count=%d&partid=%s&charge_type=%s&gm_id=%s",
		player.Id,
		params.ItemId,
		params.ChargeValue,
		params.PlatformId,
		params.ChargeType,
		c.curUser.Account,
	)
	sign := utils.String2md5(args + "fa9274fd68cf8991953b186507840e5e")
	logs.Info("sign:%v", sign)

	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	url := models.GetGameURLByNode(gameServer.Node) + "/gm_charge?" + args + "&sign=" + sign
	logs.Info("url:%v", url)
	resp, err := http.Get(url)
	c.CheckError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	c.CheckError(err)

	err = json.Unmarshal(body, &result)
	logs.Info("result:%v", string(body))
	c.CheckError(err)
	logs.Info("后台充值结果:%v", result)
	if result.Code == 0 {
		backgroundChargeLog := &models.BackgroundChargeLog{
			PlatformId:  params.PlatformId,
			ServerId:    string(player.ServerId),
			PlayerId:    params.PlayerId,
			Time:        utils.GetTimestamp(),
			ChargeType:  params.ChargeType,
			ChargeValue: params.ChargeValue,
			ItemId:        params.ItemId,
			UserId:      c.curUser.Id,
		}
		err = models.Db.Save(&backgroundChargeLog).Error
		c.CheckError(err, "写后台充值日志失败")
		c.Result(enums.CodeSuccess, "后台充值成功", 0)
	}
	c.Result(enums.CodeFail, fmt.Sprintf("后台充值失败: ErrorCode: %v Messsage", result.Code, result.Message), result.Code)
}
