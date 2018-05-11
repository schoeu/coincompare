package actions

import (
	"../utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type coinsInfo struct {
	Pid  string `json:"pid"`
	Name string `json:"name"`
}

// 获取流量信息
func GetAllFlow(c *gin.Context, infoDB *sql.DB, compareDB *sql.DB) {
	var pid, symbol string
	var ci coinsInfo
	var ciCtt []coinsInfo

	rows, err := infoDB.Query("select pid, symbol from bt_listings")
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

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   ciCtt,
	})
}
