package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 定义路由
	router.POST("/register/email", handlers.EmailRegister)
	router.POST("/register/phone", handlers.PhoneRegister)

	return router
}
