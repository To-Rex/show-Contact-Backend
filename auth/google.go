package auth

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"net/http"
	"show_contact/config"
	"show_contact/models"
	"strings"
)

func LoginGoogle(c *gin.Context) {
	var user models.User
	_ = c.BindJSON(&user)
	db := config.ConnectDB()
	//if users table does not exist, create it and insert the user
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, access_token TEXT, id_token TEXT, ids TEXT, phone TEXT, email TEXT, password TEXT, name TEXT, photo_url TEXT, status TEXT, blocked BOOLEAN, role TEXT, region TEXT, device TEXT, created_at TIMESTAMP, updated_at TIMESTAMP, token TEXT)")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating the users table" + err.Error()})
		return
	}
	var id int
	err = db.QueryRow("SELECT id FROM users WHERE email=$1 OR phone=$2", user.Email, user.Phone).Scan(&id)
	if err != nil {
		/*{
		    "access_token": "sakmdlkmowekmeocmeorckecome",
		    "id_token":"wqwqwew",
		    "ids":"sdsdswea",
		    "phone": "+998995340313",
		    "email":"dsesasweswsedase",
		    "password":"12345678",
		    "name":"xswxaewexxasdds",
		    "photo_url":"sawsasedwedwededs",
		    "status":"wsasefwefweds",
		    "blocked":false,
		    "role":"users",
		    "region":"region",
		    "device": "android",
		    "created_at":"2023-07-15 12:29:40",
		    "updated_at":"2023-07-15 12:29:40"
		}*/
		//user.Token, _ = generateToken(user.Email, user.Password, user.Role)
		//generateToken function ../config/config.go
		user.Token, _ = config.GenerateToken(user.Email, user.Password, user.Role)

		//_, err := db.Exec("INSERT INTO users (access_token, id_token, ids, phone, email, name, photo_url, status, blocked, role, region, device, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false, $9, $10, $11, $12, $13)", user.AccessToken, user.IdToken, user.Ids, user.Phone, user.Email, user.Name, user.PhotoUrl, user.Status, user.Role, user.Region, user.Device, user.CreatedAt, user.UpdatedAt) //insert the user
		_, err := db.Exec("INSERT INTO users (access_token, id_token, ids, phone, email, password, name, photo_url, blocked, role, region, device, created_at, updated_at, token) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false, $9, $10, $11, $12, $13, $14)", user.AccessToken, user.IdToken, user.Ids, user.Phone, user.Email, config.PasswordHash("root"), user.Name, user.PhotoUrl, user.Role, user.Region, user.Device, user.CreatedAt, user.UpdatedAt, user.Token) //insert the user
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating the user " + err.Error()})
			return
		}
		//return the user token and id message: User created successfully
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "token": user.Token})
		return
	} else {
		//if user is blocked, return error
		var blocked bool
		err = db.QueryRow("SELECT blocked FROM users WHERE id=$1", id).Scan(&blocked)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the user"})
			return
		}
		if blocked {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is blocked"})
			return
		}
		user.Token, _ = config.GenerateToken(user.Email, user.Password, user.Role)
		_, err := db.Exec("UPDATE users SET access_token=$1, id_token=$2, ids=$3, email=$4, name=$5, photo_url=$6, region=$7, device=$8, updated_at=$9, token=$10 WHERE id=$11", user.AccessToken, user.IdToken, user.Ids, user.Email, user.Name, user.PhotoUrl, user.Region, user.Device, user.UpdatedAt, user.Token, id) //update the user
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error updating the user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "token": user.Token})
		return
	}
}

func GetUsers(c *gin.Context) {
	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	} /*, jwt.WithAudience("foo"), jwt.WithIssuer("bar")*/)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
		return
	}
	//return all users
	db := config.ConnectDB()
	rows, err := db.Query("SELECT id, access_token, id_token, ids, phone, email, password, name, photo_url, blocked, role, region, device, created_at, updated_at FROM users")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
			return
		}
	}(rows)
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.AccessToken, &user.IdToken, &user.Ids, &user.Phone, &user.Email, &user.Password, &user.Name, &user.PhotoUrl, &user.Blocked, &user.Role, &user.Region, &user.Device, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
			return
		}
		users = append(users, user)
	}
	//return the users and the status code
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetsAllUsers(c *gin.Context) {
	db := config.ConnectDB()
	rows, err := db.Query("SELECT id, access_token, id_token, ids, phone, email, password, name, photo_url, blocked, role, region, device, created_at, updated_at FROM users")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
			return
		}
	}(rows)
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.AccessToken, &user.IdToken, &user.Ids, &user.Phone, &user.Email, &user.Password, &user.Name, &user.PhotoUrl, &user.Blocked, &user.Role, &user.Region, &user.Device, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
			return
		}
		users = append(users, user)
	}
	//return the users and the status code
	c.JSON(http.StatusOK, gin.H{"users": users})
	return
}
