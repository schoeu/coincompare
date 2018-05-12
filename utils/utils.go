package utils

import (
	"../config"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const sqlReg = "(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\\b(select|update|and|or|delete|insert|trancate|char|into|substr|ascii|declare|exec|count|master|into|drop|execute)\\b)"

// error handler.
func ErrHandle(err error) {
	if err != nil {
		log.Println(err)
	}
}

// db opener
func OpenDb(dbType string, dbStr string) *sql.DB {
	if dbType == "" {
		log.Println("No dbType")
	}
	db, err := sql.Open(dbType, dbStr)
	ErrHandle(err)

	err = db.Ping()
	ErrHandle(err)
	return db
}

func GetCurrentDate() string {
	t := time.Now().String()
	return strings.Split(t, " ")[0]
}

// check sql string
func CheckSql(s string) string {
	match, _ := regexp.Match(sqlReg, []byte(s))
	if match {
		return ""
	}
	return s
}

// {"errcode":40029,"errmsg":"invalid code, hints: [ req_id: 0SutCA0163th50 ]"}
type rawData struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errmsg     string `json:"errmsg"`
	Errcode    int    `json:"errcode"`
}

func GetUserInfoRaw(code string, encryptedData string) rawData {
	rd := rawData{}
	if code == "" {
		log.Println("No response code.")
	} else {
		url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + config.AppId + "&secret=" + config.Secret + "&js_code=" + code + "&grant_type=authorization_code"
		res, err := http.Get(url)
		ErrHandle(err)

		body, err := ioutil.ReadAll(res.Body)

		json.Unmarshal(body, &rd)
		res.Body.Close()
		ErrHandle(err)

		fmt.Println(string(body))

		//pc := NewWXBizDataCrypt(config.AppId, a.SessionKey)
		//userInfo, err := pc.Decrypt(encryptedData, iv)

	}
	return rd
}
