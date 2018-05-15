package actions

import (
	"../utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type coinsFullInfo struct {
	Count   int    `json:"count"`
	Win     int    `json:"win"`
	Lose    int    `json:"lose"`
	Name    string `json:"name"`
	Surplus string `json:"surplus"`
	Deficit string `json:"deficit"`
}

func GetList(c *gin.Context, compareDB *sql.DB) {
	var name sql.NullString
	var state sql.NullInt64
	var fullInfoArr []coinsFullInfo

	sqlStr := "select  a.symbol, b.state from bt_info.bt_listings as a left join bt_coincom.bt_coininfo as b on a.symbol = b.coin_name "

	key := c.Query("key")
	max := c.Query("max")
	pn := c.Query("pn")

	if key != "" {
		sqlStr += "where LOCATE('" + utils.CheckSql(key) + "', symbol ) > 0 "
	}
	if pn == "" {
		pn = "0"
	}
	if max == "" {
		max = "10"
	}

	pNum, err := strconv.Atoi(pn)
	maxNum, err := strconv.Atoi(max)
	utils.ErrHandle(err)
	offset := strconv.Itoa(pNum * maxNum)
	sqlStr += " limit " + max + " offset " + offset
	rows, err := compareDB.Query(sqlStr)
	utils.ErrHandle(err)

	var prevType string
	var fi coinsFullInfo
	for rows.Next() {
		err := rows.Scan(&name, &state)
		utils.ErrHandle(err)
		nameVal := name.String
		stateVal := int(state.Int64)

		if prevType != nameVal && prevType != "" {
			fullInfoArr = append(fullInfoArr, fi)
			fi = coinsFullInfo{}
		}

		if stateVal == 1 {
			fi.Win++
		} else if stateVal == 2 {
			fi.Lose++
		}
		fi.Name = nameVal
		prevType = nameVal

	}
	if fi.Name != "" {
		fullInfoArr = append(fullInfoArr, fi)
	}

	for i, v := range fullInfoArr {
		count := v.Lose + v.Win
		fullInfoArr[i].Count = count
		sVal := float64(v.Win) / float64(count)
		dVal := 1 - sVal
		if v.Lose == 0 && v.Win == 0 {
			sVal = 0.5
			dVal = 0.5
		}

		fullInfoArr[i].Surplus = fmt.Sprintf("%.1f", sVal*100) + "%"
		fullInfoArr[i].Deficit = fmt.Sprintf("%.1f", dVal*100) + "%"
	}

	err = rows.Err()
	utils.ErrHandle(err)
	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   fullInfoArr,
	})
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
			if fi.Name != "" {
				fullInfoArr = append(fullInfoArr, fi)
			}

			fi = coinsFullInfo{}
			fi.Name = name
			prevType = name
		} else {

			if state == 1 {
				fi.Win++
			} else {
				fi.Lose++
			}
		}

		fmt.Println(name, state)
	}
	if fi.Name != "" {
		fullInfoArr = append(fullInfoArr, fi)
	}

	fmt.Println(fullInfoArr)
	for i, v := range fullInfoArr {
		count := v.Lose + v.Win
		fullInfoArr[i].Count = count
		sVal := float64(v.Win) / float64(count)
		dVal := 1 - sVal
		if v.Lose == 0 && v.Win == 0 {
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
