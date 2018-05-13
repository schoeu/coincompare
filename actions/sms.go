package actions

import (
	"../config"
	"../utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/gwpp/alidayu-go"
	
	"net/http"
)

func SendSMS(c *gin.Context, _ *sql.DB, _ *sql.DB) {
	phone := c.Query("phone")
	if phone != "" {
		fmt.Println("phone", phone, config.AppKey, config.AppSecret, config.SignName, config.TemplateCode)
		client := NewTopClient(config.AppKey, config.AppSecret)
		req := request.NewAlibabaAliqinFcSmsNumSendRequest()
		req.SmsFreeSignName = config.SignName
		req.RecNum = phone
		req.SmsTemplateCode = config.TemplateCode
		req.SmsParam = ""
		response, err := client.Execute(req)
		if err != nil {
			utils.ErrHandle(err)
			return
		}
		fmt.Println("response", response)
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   "Send SMS successfully.",
		})
	}
}
