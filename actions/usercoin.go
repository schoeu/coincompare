package actions

import (
	"../utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userCoin struct {
	Count string   `json:"count"`
	Win   []string `json:"win"`
	Lose  []string `json:"lose"`
	Rate  string   `json:"rate"`
}

// 获取流量信息
func GetUserCoinInfo(c *gin.Context, _ *sql.DB, compareDB *sql.DB) {
	var pid, symbol string
	var ci coinsInfo
	var ciCtt []coinsInfo

	code := c.PostForm("code")
	iv := c.PostForm("iv")
	cryptData := c.PostForm("cryptData")
	if code != "" {
		userRawInfo := utils.GetUserInfoRaw(code, cryptData, iv)
		fmt.Println("userRawInfo",userRawInfo.UnionId)

		rows, err := compareDB.Query("SELECT coin_name, state FROM bt_coincom.bt_coininfo where uid = ?", "3")
		utils.ErrHandle(err)

		for rows.Next() {
			err := rows.Scan(&pid, &symbol)
			utils.ErrHandle(err)

			ci.Pid = pid
			ci.Name = symbol

			ciCtt = append(ciCtt, ci)
		}
		err = rows.Err()
		utils.ErrHandle(err)
		defer rows.Close()
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   ciCtt,
	})
}
