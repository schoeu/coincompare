package actions

import (
	"../utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

// init user info.
func InitUser(c *gin.Context, compareDB *sql.DB) {
	var hasState string
	var uc UserCoin

	uid := c.PostForm("uid")
	nickName := ""
	if uid == "" {
		code := c.PostForm("code")
		iv := c.PostForm("iv")
		cryptData := c.PostForm("cryptData")
		userRawInfo := utils.GetUserInfoRaw(code, cryptData, iv)
		nickName = userRawInfo.NickName
		uid = userRawInfo.UnionId
	}

	state := c.PostForm("state")
	coinName := c.PostForm("coinname")
	if uid != "" {
		_, err := compareDB.Exec("INSERT INTO bt_user(uid, name) VALUE(?, ?) ON DUPLICATE KEY UPDATE name = ?", uid, nickName, uid)
		utils.ErrHandle(err)

		rows, err := compareDB.Query("SELECT state FROM bt_coininfo where uid = ? and coin_name = ?", uid, coinName)
		utils.ErrHandle(err)
		for rows.Next() {
			err := rows.Scan(&hasState)
			utils.ErrHandle(err)
		}
		err = rows.Err()
		utils.ErrHandle(err)
		defer rows.Close()

		if hasState == "" {
			_, err := compareDB.Exec("insert into bt_coininfo (uid, coin_name, state) values (?, ?, ?)", uid, coinName, state)
			utils.ErrHandle(err)
		} else {
			_, err := compareDB.Exec("update bt_coininfo set coin_name = ?, state = ? where uid = ?", coinName, state, uid)
			utils.ErrHandle(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   uc,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "No uid for request.",
			"data":   "",
		})
	}

}
