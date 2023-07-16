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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Next()
	})
	r.POST("/auth/login", auth.Login)
	r.GET("/auth/getdata", auth.GetUsers)
	err := r.Run()
	if err != nil {
		return
	}
}
