package main

import (
	"./config"
	"./routers"
	"./utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	infoDb := utils.OpenDb("mysql", config.InfoDB)
	compareDb := utils.OpenDb("mysql", config.CompareDB)

	defer infoDb.Close()
	defer compareDb.Close()

	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is running.")
	})

	// API路由处理
	apiRouters(app, infoDb, compareDb)
	app.Run(config.Port)
}

// API路由处理
func apiRouters(router *gin.Engine, infoDb *sql.DB, compareDb *sql.DB) {
	apis := router.Group("/api")

	apis.GET("/:type", func(c *gin.Context) {
		dataType := c.Param("type")
		actions := routers.RouterMap[dataType]
		if actions != nil {
			actions(c, infoDb, compareDb)
		}
	})
}