package url

const (
	GAME_URL = "/game" //游戏主URL
	API_URL  = "/api"  //后台接口URL
	DOTA_URL = "/dot"  //打点URL
	//------------------------------------------------------------------------------------------------------------------
	//玩家信息模块
	USER_URL        = "/user"
	USER_LOGIN_URL  = GAME_URL + USER_URL + "/login"
	USER_INFO_URL   = GAME_URL + USER_URL + "/info"
	USER_LEVEL_URL  = GAME_URL + USER_URL + "/level"
	USER_UPDATE_URL = GAME_URL + USER_URL + "/update"
)
