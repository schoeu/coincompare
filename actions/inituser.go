package actions

import (
	"../utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

// init user info.
func InitUser(c *gin.Context, _ *sql.DB, compareDB *sql.DB) {
	var suid string
	var uc UserCoin

	code := c.PostForm("code")
	iv := c.PostForm("iv")
	cryptData := c.PostForm("cryptData")
	state := c.PostForm("state")
	coinName := c.PostForm("coinname")
	if code != "" {
		userRawInfo := utils.GetUserInfoRaw(code, cryptData, iv)
		uid := userRawInfo.UnionId
		_, err := compareDB.Exec("INSERT INTO bt_user(uid, name) VALUE(?, ?) ON DUPLICATE KEY UPDATE name = ?", uid, userRawInfo.NickName, uid)
		utils.ErrHandle(err)

		err = compareDB.QueryRow("SELECT uid FROM bt_coininfo where uid = ? and state = ?", uid, state).Scan(&suid)
		utils.ErrHandle(err)
		if suid == "" {
			_, err := compareDB.Exec("insert into bt_coininfo (uid, coin_name, state) values (?, ?, ?)", uid, coinName, state)
			utils.ErrHandle(err)
		} else {
			_, err := compareDB.Exec("update bt_coininfo set coin_name = ?, state = ? where uid = ?", coinName, state, uid)
			utils.ErrHandle(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   uc,
	})
}
