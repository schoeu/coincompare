package routers

import (
	"../actions"
	"database/sql"
	"github.com/gin-gonic/gin"
)

var RouterMap = map[string]func(*gin.Context, *sql.DB, *sql.DB){
	"coinslist": actions.GetAllFlow,
}
