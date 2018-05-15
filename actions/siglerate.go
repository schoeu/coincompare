package actions

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

// get rate of all users.
func GetSigleRate(c *gin.Context, compareDB *sql.DB) {
	name := c.Query("name")
	param := []string{"'" + name + "'"}
	if name != "" {
		infos := getoinInfo(param, compareDB)
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   infos,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "No name for request.",
			"data":   "",
		})
	}
}
