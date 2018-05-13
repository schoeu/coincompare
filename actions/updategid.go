package actions

import (
	"../utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

// update share info.
func UpdteGid(c *gin.Context, _ *sql.DB, compareDB *sql.DB) {
	code := c.PostForm("code")
	iv := c.PostForm("iv")
	cryptData := c.PostForm("cryptData")
	uid := c.PostForm("uid")

	if code != "" {
		userRawInfo := utils.GetUserInfoRaw(code, cryptData, iv)
		gid := userRawInfo.OpenGId
		if uid != "" && gid != "" {
			_, err := compareDB.Exec("INSERT INTO bt_group(group_id, uid) VALUE(?, ?) ", gid, uid)
			utils.ErrHandle(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   gid,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "No code to request.",
			"data":   "",
		})
	}
}
