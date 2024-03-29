package models

import (
	"errors"
	"fmt"
	"github.com/chnzrb/myadmin/utils"
	"strings"
)

type PropConsumeStatistics struct {
	OpType int     `json:"opType"`
	Count  int     `json:"count"`
	Rate   float32 `json:"rate"`
}
type PropConsumeStatisticsQueryParam struct {
	PlayerName string
	PlayerId   int
	PlatformId string
	ServerId       string `json:"serverId"`
	ChannelList [] string `json:"channelList"`
	StartTime  int
	EndTime    int
	PropType   int
	PropId     int
	Type       int
}

func GetPropConsumeStatistics(param *PropConsumeStatisticsQueryParam) ([]*PropConsumeStatistics, error) {
	if param.PropType == 0 || param.PropId == 0 {
		return nil, errors.New("请选择道具")
	}
	gameServer, err := GetGameServerOne(param.PlatformId, param.ServerId)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	list := make([]*PropConsumeStatistics, 0)

	whereArray := make([] string, 0)
	whereArray = append(whereArray, fmt.Sprintf("type = %d", param.Type))
	whereArray = append(whereArray, fmt.Sprintf("prop_type = %d", param.PropType))
	whereArray = append(whereArray, fmt.Sprintf("prop_id = %d", param.PropId))
	//whereArray = append(whereArray, fmt.Sprintf("server_id = %s", param.ServerId))
	whereArray = append(whereArray, fmt.Sprintf("player_id in (select id from player where channel in(%s) and server_id = '%s')", GetSQLWhereParam(param.ChannelList),  param.ServerId))
	if param.PlayerId > 0 {
		whereArray = append(whereArray, fmt.Sprintf("player_id = %d", param.PlayerId))
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}


	order := ""

	if param.Type == 0 {
		order = "order by count desc"
	} else {
		order = "order by count asc"
	}

	sql := fmt.Sprintf(
		` select log_type as op_type, sum(value) as count from consume_statistics  %s group by log_type %s; `, whereParam, order)
	err = gameDb.Raw(sql).Scan(&list).Error
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	var sum = 0
	for _, e := range list {
		sum += e.Count
	}
	for _, e := range list {
		e.Rate = float32(e.Count) / float32(sum) * 100
	}
	return list, nil
}

//
//func GetPropConsumeStatistics(param *PropConsumeStatisticsQueryParam) ([]*PropConsumeStatistics, error) {
//	if param.PropType == 0 || param.PropId == 0 {
//		return nil, errors.New("请选择道具")
//	}
//	gameServer, err := GetGameServerOne(param.PlatformId, param.ServerId)
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	node := gameServer.Node
//	gameDb, err := GetGameDbByNode(node)
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	defer gameDb.Close()
//	list := make([]*PropConsumeStatistics, 0)
//
//	var selectPlayer string
//	if param.PlayerId > 0 {
//		selectPlayer = fmt.Sprintf("and player_id = %d", param.PlayerId)
//	}
//
//	order := ""
//
//	if param.Type == 0 {
//		order = "order by count desc"
//	} else {
//		order = "order by count asc"
//	}
//
//	sql := fmt.Sprintf(
//		` select log_type as op_type, sum(value) as count from consume_statistics where type = %d and prop_type = %d and prop_id = %d %s group by log_type %s; `, param.Type, param.PropType, param.PropId, selectPlayer, order)
//	err = gameDb.Raw(sql).Scan(&list).Error
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	var sum = 0
//	for _, e := range list {
//		sum += e.Count
//	}
//	for _, e := range list {
//		e.Rate = float32(e.Count) / float32(sum) * 100
//	}
//	return list, nil
//}


//func GetPropConsumeStatistics(param *PropConsumeStatisticsQueryParam) ([]*PropConsumeStatistics, error) {
//	if param.PropType == 0 || param.PropId == 0 {
//		return nil, errors.New("请选择道具")
//	}
//	gameDb, err := GetGameDbByNode(param.Node)
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	defer gameDb.Close()
//	list := make([]*PropConsumeStatistics, 0)
//	var changeValue string
//	if param.Type == 0 {
//		changeValue = "change_value < 0"
//	} else {
//		changeValue = "change_value > 0"
//	}
//	var timeRange string
//	if param.StartTime > 0 {
//		timeRange = fmt.Sprintf("and op_time between %d and %d", param.StartTime, param.EndTime)
//	}
//
//	var selectPlayer string
//	if param.PlayerId > 0 {
//		timeRange = fmt.Sprintf("and player_id = %d", param.PlayerId)
//	}
//
//	sql := fmt.Sprintf(
//		` select op_type, sum(change_value) as count from player_prop_log where %s and prop_type = ? and prop_id = ? %s %s group by op_type; `, changeValue, timeRange, selectPlayer)
//	err = gameDb.Raw(sql, param.PropType, param.PropId).Scan(&list).Error
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	var sum = 0
//	for _, e := range list {
//		sum += e.Count
//	}
//	for _, e := range list {
//		e.Rate = float32(e.Count) / float32(sum) * 100
//	}
//	return list, nil
//}
