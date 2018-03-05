package enums

type JsonResultCode int
const (
	ChosePlatformId = "chose_platform_id"
	ChoseServerId = "chose_server_id"
)
const (
	JRCodeSucc JsonResultCode = iota
	JRCodeFailed
	CodeSuccess = 0
	CodeFail = 1
	CodeUnauthorized = 401 //未授权访问
	CodeNoLogin = 50014 //未登录
)
const (
	MSG_OK  = 0
	MSG_ERR = -1
)

const (
	Deleted = iota - 1
	Disabled
	Enabled
)
