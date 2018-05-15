package actions

import (
	"../utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type signUp struct {
	IsSignUp bool `json:"isSignUp"`
}

// get signup info.
func IsSignUp(c *gin.Context, compareDB *sql.DB) {
	var state int

	uid := c.PostForm("uid")
	if uid == "" {
		code := c.PostForm("code")
		iv := c.PostForm("iv")
		cryptData := c.PostForm("cryptData")
		userRawInfo := utils.GetUserInfoRaw(code, cryptData, iv)
		uid = userRawInfo.UnionId

		fmt.Println(userRawInfo)
	}

	if uid != "" {
		rows, err := compareDB.Query("SELECT EXISTS(select * from bt_user where uid = ?)", uid)
		utils.ErrHandle(err)
		for rows.Next() {
			err := rows.Scan(&state)
			utils.ErrHandle(err)
			fmt.Println(state, uid)
		}
		err = rows.Err()
		utils.ErrHandle(err)
		defer rows.Close()
	}

	su := signUp{}
	if state == 1 {
		su.IsSignUp = true
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   su,
	})
}
