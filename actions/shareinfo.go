package actions

import (
	"../utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// get signup info.
func ShareInfo(c *gin.Context, _ *sql.DB, compareDB *sql.DB) {
	code := c.PostForm("code")
	iv := c.PostForm("iv")
	cryptData := c.PostForm("cryptData")

	if code != "" {
		userRawInfo := utils.GetUserInfoRaw(code, cryptData, iv)
		uid := userRawInfo.UnionId
		if uid != "" {
			fmt.Println("userRawInfo~", userRawInfo)

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   "su",
	})
}
