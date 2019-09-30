package utils

//提示常量
const LOGIN_ERROR_NOTICE = "账号/密码错误"

//token验证常量
const TOKEN_PARAM_REQUIRE = "token参数缺失"
const TOKEN_INVDLID = "token验证无效"
const TOKEN_PARSE_ERROR = "token解析错误"

//验证常量
const TOKEN_TYPE = "Bearer"
const TOKEN_TYPE_ERROR = "token加密类型错误"

//加密常量
const TOKEN_SECRET = "196ff70efaf6913f"

//Api Prefix
const API_PREFIX = "/api"

//Server Port
const SERVER_PORT = "8086"

//LogPath
const LOG_PATH = "D:/gopath/src/edu_api/log"

//Default log name
const LOG_NAME = "request"

//redis缓存key
const (
	LATEST_MEDIUM_PLAY              = "latest:medium:play:user:%d"
	LATEST_LESION_INFO              = "latest:lesion:info"
	LATEST_CLASS_WATCH_INFO         = "latest:class:%d:watch:info"
	LATEST_CLASS_CHAPTER_WATCH_INFO = "latest:class:%d:chapter:%d:watch:info"
)
