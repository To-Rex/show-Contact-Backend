package services

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"show_contact/auth"
)

func StartServices() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	})
	r.POST("/auth/login", auth.Login)
	r.GET("/auth/getdata", auth.GetUsers)
	err := r.Run()
	if err != nil {
		return
	}
}
