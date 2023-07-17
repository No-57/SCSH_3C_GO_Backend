package main

import (
	"SCSH_3C_GO_Backend/db"
	"SCSH_3C_GO_Backend/pkg/Utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
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

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random"
)

const (
	GoogleClientID     = ""
	GoogleClientSecret = ""
	LineClientID       = ""
	LineClientSecret   = ""
)

var (
	LineRedirectURL   = "http://localhost:8080/login/line/callback"
	GoogleRedirectURL = "http://localhost:8080/callback"
)

func main() {

	Utils.InitProperties("conf/", "config", "properties")
	//Utils.InitProperties("/Users/kaoweicheng/GolandProjects/NO57_backend/conf/", "config", "properties")

	//db初始化
	db.InitDB()
	engine := db.Engine

	r := gin.Default()

	r.LoadHTMLGlob("templates/html/*")

	r.GET("/", handleHome)
	//google登入
	r.GET("/login/google", handleGoogleLogin)
	r.GET("/callback", handleGoogleCallback)

	//line登入
	r.GET("/login/line", handleLineLogin)
	r.GET("/login/line/callback", handleLineCallback)

	//app自己的登入
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

func handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func handleGoogleLogin(c *gin.Context) {
	// google oauth config
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  GoogleRedirectURL,
		ClientID:     GoogleClientID,
		ClientSecret: GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleGoogleCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Fatal(err)
	}

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	var userInfo struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Picture  string `json:"picture"`
		Verified bool   `json:"verified_email"`
	}

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		log.Fatal(err)
	}

	// 基本資料
	fmt.Println("User ID:", userInfo.ID)
	fmt.Println("Email:", userInfo.Email)
	fmt.Println("Name:", userInfo.Name)
	fmt.Println("Picture:", userInfo.Picture)
	fmt.Println("Verified Email:", userInfo.Verified)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func handleLineLogin(c *gin.Context) {
	// 创建OAuth2配置
	config := &oauth2.Config{
		ClientID:     LineClientID,
		ClientSecret: LineClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
			TokenURL: "https://api.line.me/oauth2/v2.1/token",
		},
		RedirectURL: LineRedirectURL,
		Scopes:      []string{"profile", "openid"},
	}

	// 生成認證URL
	authURL := config.AuthCodeURL("state", oauth2.AccessTypeOnline)

	// 重定向到認證URL
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func handleLineCallback(c *gin.Context) {
	// 獲取授權碼
	code := c.Query("code")

	// 創建OAuth2配置
	config := &oauth2.Config{
		ClientID:     LineClientID,
		ClientSecret: LineClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
			TokenURL: "https://api.line.me/oauth2/v2.1/token",
		},
		RedirectURL: LineRedirectURL,
		Scopes:      []string{"profile", "openid"},
	}

	// 通過授權碼獲取token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatal(err)
	}

	// 使用令牌从Line获取用户信息
	req, err := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 取得回傳資訊
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 印出user data
	fmt.Println(string(body))

	// 返回指定頁面
	c.HTML(http.StatusOK, "index.html", nil)
}
