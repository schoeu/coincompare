package routers

import (
	"../actions"
	"database/sql"
	"github.com/gin-gonic/gin"
)

var GETRouterMap = map[string]func(*gin.Context, *sql.DB, *sql.DB){
	"coinslist": actions.GetAllFlow,
}
var POSTRouterMap = map[string]func(*gin.Context, *sql.DB, *sql.DB){
	"usercoininfo": actions.GetUserCoinInfo,
}
