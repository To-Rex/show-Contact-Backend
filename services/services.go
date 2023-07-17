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
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Content-Type", "application/json")

		c.Next()
	})

	r.POST("/auth/login", auth.Login)
	r.GET("/auth/getdata", auth.GetUsers)
	r.OPTIONS("/auth/getdata", auth.GetUsers)
	r.GET("/auth/getdata1", auth.GetsAllUsers)
	err := r.Run()
	if err != nil {
		return
	}
}
