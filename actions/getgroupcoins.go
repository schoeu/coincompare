package actions

import (
	"../utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type groupCoinsInfo struct {
	Win   int    `json:"win"`
	Lose  int    `json:"lose"`
	WRate string `json:"wRate"`
	LRate string `json:"lRate"`
}

// get coin info in a wx group.
func GetGroupCoins(c *gin.Context, _ *sql.DB, compareDB *sql.DB) {
	code := c.PostForm("code")
	iv := c.PostForm("iv")
	cryptData := c.PostForm("cryptData")
	name := c.PostForm("name")
	gci := groupCoinsInfo{}
	if code != "" {
		var state int
		userRawInfo := utils.GetUserInfoRaw(code, cryptData, iv)
		gid := userRawInfo.OpenGId
		if name != "" && gid != "" {
			rows, err := compareDB.Query("select state from bt_coininfo where uid in ( SELECT distinct uid FROM bt_group where group_id = ? ) and coin_name = ?", gid, name)
			utils.ErrHandle(err)
			for rows.Next() {
				err := rows.Scan(&state)
				utils.ErrHandle(err)

				if state == 1 {
					gci.Win++
				} else {
					gci.Lose++
				}
			}
			err = rows.Err()
			utils.ErrHandle(err)
			defer rows.Close()

			count := gci.Win + gci.Lose
			if count != 0 {
				wVal := float64(gci.Win) / float64(count)
				lVal := 1 - wVal
				gci.WRate = fmt.Sprintf("%.1f", wVal*100) + "%"
				gci.LRate = fmt.Sprintf("%.1f", lVal*100) + "%"
			} else {
				gci.WRate = "50%"
				gci.LRate = "50%"
			}

		}
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   gci,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "No code to request.",
			"data":   "",
		})
	}
}
