package main

import (
	"fmt"
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

	compareDb := utils.OpenDb("mysql", config.CompareDB)

	defer compareDb.Close()

	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is running.")
	})

	// API路由处理
	apiRouters(app, compareDb)
	app.Run(config.Port)
}

// API路由处理
func apiRouters(router *gin.Engine, compareDb *sql.DB) {
	apis := router.Group("/api")

	// get method actions
	apis.GET("/:type", func(c *gin.Context) {
		fmt.Println("URL", c.Request.URL)
		dataType := c.Param("type")
		actions := routers.GETRouterMap[dataType]
		if actions != nil {
			actions(c, compareDb)
		}
	})
	// post method actions
	apis.POST("/:type", func(c *gin.Context) {
		dataType := c.Param("type")
		actions := routers.POSTRouterMap[dataType]
		if actions != nil {
			actions(c, compareDb)
		}
	})
}
