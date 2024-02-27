package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"your_project/database" // 替换为你的项目路径
	"your_project/models"   // 替换为你的项目路径
)

// EmailRegister 处理邮箱注册请求
func EmailRegister(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 这里应加入密码加密逻辑

	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user": user})
}
