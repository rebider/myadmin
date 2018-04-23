package models

import (
	"errors"
	"fmt"
	"github.com/chnzrb/myadmin/utils"
)

type PropConsumeStatistics struct {
	OpType int     `json:"opType"`
	Count  int     `json:"count"`
	Rate   float32 `json:"rate"`
}
type PropConsumeStatisticsQueryParam struct {
	PlayerName string
	PlayerId   int
	PlatformId int
	Node       string `json:"serverId"`
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
	gameDb, err := GetGameDbByNode(param.Node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	list := make([]*PropConsumeStatistics, 0)
	var changeValue string
	if param.Type == 0 {
		changeValue = "change_value < 0"
	} else {
		changeValue = "change_value > 0"
	}
	var timeRange string
	if param.StartTime > 0 {
		timeRange = fmt.Sprintf("and op_time between %d and %d", param.StartTime, param.EndTime)
	}

	var selectPlayer string
	if param.PlayerId > 0 {
		timeRange = fmt.Sprintf("and player_id = %d", param.PlayerId)
	}

	sql := fmt.Sprintf(
		` select op_type, sum(change_value) as count from player_prop_log where %s and prop_type = ? and prop_id = ? %s %s group by op_type; `, changeValue, timeRange, selectPlayer)
	err = gameDb.Raw(sql, param.PropType, param.PropId).Scan(&list).Error
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
