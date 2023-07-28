package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"show_contact/config"
	"show_contact/models"
	"strings"
	"time"
)

func GetUser(c *gin.Context) {
	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if claims["created_at"].(float64) < float64(time.Now().Unix()) {
			return nil, nil
		}
		return []byte("secret"), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users or token is expired"})
		return
	}
	db := config.ConnectDB()
	var user models.User
	err = db.QueryRow("SELECT access_token, id_token, ids, phone, email, password, name, photo_url, blocked, role, region, device, created_at, updated_at FROM users WHERE email = $1", claims["email"]).Scan(&user.AccessToken, &user.IdToken, &user.Ids, &user.Phone, &user.Email, &user.Password, &user.Name, &user.PhotoUrl, &user.Blocked, &user.Role, &user.Region, &user.Device, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
