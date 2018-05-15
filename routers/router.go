package routers

import (
	"../actions"
	"database/sql"
	"github.com/gin-gonic/gin"
)

var GETRouterMap = map[string]func(*gin.Context, *sql.DB){
	"coinslist":     actions.GetList,
	"siglerate":     actions.GetSigleRate,
	"getgroupcoins": actions.GetGroupCoins,
}
var POSTRouterMap = map[string]func(*gin.Context, *sql.DB){
	"usercoininfo": actions.GetUserCoinInfo,
	"inituser":     actions.InitUser,
	"issignup":     actions.IsSignUp,
	"allrate":      actions.GetList,
	"login":        actions.Login,
	"updategid":    actions.UpdteGid,
	"getuid":       actions.GetUid,
}
