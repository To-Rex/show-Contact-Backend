package services

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"show_contact/auth"
)

func StartServices() {
	r := gin.Default()
	r.POST("/auth/login", auth.Login)
	r.GET("/auth/getdata", auth.GetUsers)
	err := r.Run()
	if err != nil {
		return
	}
}
