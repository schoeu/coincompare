package actions

import (
	"../utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type coinsInfo struct {
	Pid  string `json:"pid"`
	Name string `json:"name"`
}

type coinsFullInfo struct {
	Count   int    `json:"count"`
	win     int    `json:"win"`
	lose    int    `json:"lose"`
	Name    string `json:"name"`
	Surplus string `json:"surplus"`
	Deficit string `json:"deficit"`
}

// get coin list data.
func GetList(c *gin.Context, infoDB *sql.DB, compareDB *sql.DB) {
	var pid, symbol string
	var cList []string

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

	rows, err := infoDB.Query(sqlStr)
	utils.ErrHandle(err)

	for rows.Next() {
		err := rows.Scan(&pid, &symbol)
		utils.ErrHandle(err)

		cList = append(cList, "'"+symbol+"'")
	}
	err = rows.Err()
	utils.ErrHandle(err)
	defer rows.Close()

	if len(cList) > 0 {
		v := getoinInfo(cList, compareDB)
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   v,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "Coins list is empty.",
			"data":   "",
		})
	}
}

func getoinInfo(coins []string, db *sql.DB) []coinsFullInfo {
	var name string
	var state int
	var fullInfoArr []coinsFullInfo

	sqlStr := "select coin_name, state from bt_coininfo where coin_name in (" + strings.Join(coins, ",") + ") order by coin_name "
	fmt.Println("sqlStr", sqlStr)
	rows, err := db.Query(sqlStr)
	utils.ErrHandle(err)

	var prevType string
	var fi coinsFullInfo
	for rows.Next() {
		err := rows.Scan(&name, &state)
		utils.ErrHandle(err)

		if prevType != name {
			fullInfoArr = append(fullInfoArr, fi)
			fi = coinsFullInfo{}
			fi.Name = name
			prevType = name
		} else {

			if state == 1 {
				fi.win++
			} else {
				fi.lose++
			}
		}

		fmt.Println(name, state)
	}

	fmt.Println(fullInfoArr)
	for i, v := range fullInfoArr {
		count := v.lose + v.win
		fullInfoArr[i].Count = count
		sVal := float64(v.win) / float64(count)
		dVal := 1 - sVal
		if v.lose == 0 && v.win == 0 {
			sVal = 0.5
			dVal = 0.5
		}

		fullInfoArr[i].Surplus = fmt.Sprintf("%.1f", sVal*100) + "%"
		fullInfoArr[i].Deficit = fmt.Sprintf("%.1f", dVal*100) + "%"
	}

	err = rows.Err()
	utils.ErrHandle(err)
	defer rows.Close()

	return fullInfoArr
}
