package actions

import (
	"../utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type coinsInfo struct {
	Pid  string `json:"pid"`
	Name string `json:"name"`
}

// get coin list data.
func GetList(c *gin.Context, infoDB *sql.DB, compareDB *sql.DB) {
	var pid, symbol string
	var ci coinsInfo
	var ciCtt []coinsInfo

	sqlStr := "select pid, symbol from bt_listings "
	key := c.Query("key")
	max := c.Query("max")
	offset := c.Query("offset")

	if key != "" {
		sqlStr += "where LOCATE('" + utils.CheckSql(key) + "', symbol ) > 0 "
	}
	if offset == "" {
		offset = "0"
	}
	if max == "" {
		max = "10"
	}

	sqlStr += " limit " + max + " offset " + offset

	fmt.Println(sqlStr)

	rows, err := infoDB.Query(sqlStr)
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
