package models

import "github.com/chnzrb/myadmin/enums"

type JsonResult struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Obj  interface{}          `json:"obj"`
}

type Result struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Data  interface{}          `json:"data"`
}

type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int    `json:"limit"`
}
