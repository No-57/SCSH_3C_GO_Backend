package main

import (
	"NO57_backend/db"
	"NO57_backend/pkg/Utils"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

// USER_DATA table 的結構
type userData struct {
	UserID   string `xorm:"'USER_ID'"`
	UserName string `xorm:"'USER_NAME'"`
	UserPwd  string `xorm:"'USER_PWD'"`
	Sex      string `xorm:"'SEX'"`
	Nickname string `xorm:"'NICKNAME'"`
	Email    string `xorm:"'EMAIL'"`
}

func main() {

	Utils.InitProperties("conf/", "config", "properties")
	//Utils.InitProperties("/Users/kaoweicheng/GolandProjects/NO57_backend/conf/", "config", "properties")

	//db初始化
	db.InitDB()
	engine := db.Engine

	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 收到的請求帳密
		fmt.Println("usernamd=", request.Username)
		fmt.Println("password=", request.Password)

		// 測試資料庫連線
		if err := engine.Ping(); err != nil {
			log.Fatal(err)
		}

		// 根據用戶名和密碼查詢用戶
		user := userData{}
		exists, err := engine.Table("USER_DATA").Where("USER_ID = ? AND USER_PWD = ?", request.Username, request.Password).Get(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫錯誤"})
			return
		}

		if exists {
			// 用戶存在且密碼正確
			c.JSON(http.StatusOK, gin.H{"message": "登錄成功"})
		} else {
			// 用戶不存在或密碼不正確
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的用戶名稱或密碼"})
		}
	})

	// port
	r.Run(":8080")

}
