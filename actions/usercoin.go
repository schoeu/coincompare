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

	code := c.Query("code")
	if code != "" {
		userRawInfo := utils.GetUserInfoRaw(code, "")
		fmt.Println(userRawInfo)

		rows, err := compareDB.Query("SELECT coin_name, state FROM bt_coincom.bt_coininfo where uid = ?", userRawInfo.Unionid)
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
