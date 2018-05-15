package actions

import (
	"../utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context, _ *sql.DB, compareDB *sql.DB) {
	uid := c.PostForm("uid")
	if uid == "" {
		code := c.PostForm("code")
		iv := c.PostForm("iv")
		cryptData := c.PostForm("cryptData")
		userRawInfo := utils.GetUserInfoRaw(code, cryptData, iv)
		uid = userRawInfo.UnionId
	}

	phone := c.PostForm("phone")
	wallet := c.PostForm("wallet")

	if uid != "" {
		_, err := compareDB.Query("INSERT INTO bt_user(uid, phone, wallet) VALUE(?, ?, ?) ON DUPLICATE KEY UPDATE phone = ?, wallet = ?", uid, phone, wallet, phone, wallet)
		utils.ErrHandle(err)
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "No uid for update.",
			"data":   "",
		})
	}

}
