package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"show_contact/config"
	"show_contact/models"
)

func Login(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users  " + err.Error()})
		return
	}
	db := config.ConnectDB()
	//chesk user.email or user.phone and user.password if db in users table have this user return user.token
	//err = db.QueryRow("SELECT phone, email, password, FROM users WHERE email = $1 OR phone = $1 AND password = $2", user.Email, user.Password).Scan(&user.Phone, &user.Email, &user.Password)
	err = db.QueryRow("SELECT id, access_token, id_token, ids, phone, email, password, name, photo_url, blocked, role, region, device, created_at, updated_at FROM users WHERE email = $1 OR phone = $1 AND password = $2", user.Email, user.Password).Scan(&user.Id, &user.AccessToken, &user.IdToken, &user.Ids, &user.Phone, &user.Email, &user.Password, &user.Name, &user.PhotoUrl, &user.Blocked, &user.Role, &user.Region, &user.Device, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
		return
	}
	//return the user token and the status code
	token := user.Token

	if true {
		token, err = config.GenerateToken(user.Email, user.Name, user.Role)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
			return
		}
		_, err = db.Exec("UPDATE users SET token = $1 WHERE email = $2", token, user.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
	return
}
