package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var balance = 0

type Result struct {
	Data struct {
		Amount  int    `json:"amount"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"data"`
}

var result = Result{}

func main() {
	router := gin.Default()
	router.GET("/deposit/:input", deposit)
	router.GET("/withdraw/:input", withdraw)
	router.GET("/balance/", getBalance)

	router.Run(":80")
}

// getBalance 取得帳戶內餘額
func getBalance(context *gin.Context) {
	result.Data.Amount = balance // 返回的Amount為當前餘額
	result.Data.Status = "ok"    // 查詢時，可將預設狀態設為成功
	result.Data.Message = ""     // 成功時Message提示為空
	context.JSON(http.StatusNotFound, result)
	//context.JSON(http.StatusOK, result)
}

// deposit 儲值、存款
func deposit(context *gin.Context) {
	input := context.Param("input")
	amount, err := strconv.Atoi(input)

	result.Data.Status = "failed" // 存款操作時，可將預設狀態設為失敗
	result.Data.Message = ""

	if err == nil {
		if amount <= 0 {
			result.Data.Amount = 0 // 操作未成功，返回金額為0
			result.Data.Message = "操作失敗，存款金額需大於0元！"
		} else {
			balance += amount
			result.Data.Amount = balance // 操作成功，返回的Amount為儲值後的餘額
			result.Data.Status = "ok"    // 操作成功
		}
	} else {
		result.Data.Amount = 0 // 操作未成功，返回金額為0
		result.Data.Message = "操作失敗，輸入有誤！"
	}
	context.JSON(http.StatusOK, result)
}

// withdraw 提款
func withdraw(context *gin.Context) {
	result.Data.Status = "failed" // 提款操作時，可將預設狀態設為失敗
	result.Data.Message = ""

	input := context.Param("input")
	amount, err := strconv.Atoi(input)

	if err == nil {
		if amount <= 0 {
			result.Data.Amount = 0 // 操作未成功，返回金額為0
			result.Data.Message = "操作失敗，提款金額需大於0元！"
		} else {
			if balance-amount < 0 {
				result.Data.Amount = 0 // 操作未成功，返回金額為0
				result.Data.Message = "操作失敗，餘額不足！"
			} else {
				balance -= amount
				result.Data.Amount = balance // 操作成功，返回的Amount為提款後的餘額
				result.Data.Status = "ok"
			}
		}
	} else {
		result.Data.Amount = 0 // 操作未成功，返回金額為0
		result.Data.Message = "操作失敗，輸入有誤！"
	}
	context.JSON(http.StatusOK, result)
}
