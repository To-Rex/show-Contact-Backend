package services

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"show_contact/auth"
	"show_contact/contacts"
)

type CORSRouterDecorator struct {
	R http.Handler
}

func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	if req.Method == "OPTIONS" {
		return
	}
	c.R.ServeHTTP(rw, req)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Token, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func StartServices() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/auth/login/google", auth.LoginGoogle)
	r.POST("/auth/login", auth.Login)
	r.GET("/auth/getdata", auth.GetUsers)
	r.GET("/auth/getdata1", auth.GetsAllUsers)
	r.POST("/data/contact/import", ImPortContact)
	r.POST("/data/contact/export", ExportContact)
	r.POST("/contact", contacts.AddContact)
	//search contact
	r.GET("/contact/search", contacts.SearchContact)
	err := r.Run()
	if err != nil {
		return
	}
}
