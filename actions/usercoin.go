package actions

import (
	"../utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserCoin struct {
	Count int      `json:"count"`
	Win   []string `json:"win"`
	Lose  []string `json:"lose"`
	Rate  string   `json:"rate"`
}

// 获取流量信息
func GetUserCoinInfo(c *gin.Context, _ *sql.DB, compareDB *sql.DB) {
	var name string
	var state int
	var uc UserCoin

	code := c.PostForm("code")
	iv := c.PostForm("iv")
	cryptData := c.PostForm("cryptData")
	if code != "" {
		userRawInfo := utils.GetUserInfoRaw(code, cryptData, iv)
		fmt.Println(userRawInfo)
		// rows, err := compareDB.Query("SELECT coin_name, state FROM bt_coincom.bt_coininfo where uid = ?", userRawInfo.UnionId)
		rows, err := compareDB.Query("SELECT coin_name, state FROM bt_coincom.bt_coininfo where uid = ?", "oH5VJ1MIfLhswr0LFYRf5RrAEcsQ")
		utils.ErrHandle(err)
		count := 0
		for rows.Next() {
			err := rows.Scan(&name, &state)
			utils.ErrHandle(err)

			if state == 0 {
				uc.Lose = append(uc.Lose, name)
			} else {
				uc.Win = append(uc.Win, name)
			}
			count++
		}
		uc.Count = count
		uc.Rate = fmt.Sprintf("%.1f", float64(len(uc.Win))/float64(count)*100)

		err = rows.Err()
		utils.ErrHandle(err)
		defer rows.Close()
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   uc,
	})
}
