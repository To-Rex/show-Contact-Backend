package auth

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

type User struct {
	Id          int    `json:"id"`
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
	Ids         string `json:"ids"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	PhotoUrl    string `json:"photo_url"`
	Blocked     bool   `json:"blocked"`
	Role        string `json:"role"`
	Region      string `json:"region"`
	Device      string `json:"device"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

const (
	host     = "containers-us-west-73.railway.app" //host
	port     = 6725                                //port
	user     = "postgres"                          //foydalanuvchi
	password = "RHnJjy7DZt2HGGw7kn28"              //parol
	dbname   = "railway"                           //baza nomi
)

func connectDB() *sql.DB { //Dastur bilan bazaga ulanish
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname) //bazaga ulanish uchun ma'lumotlar
	db, err := sql.Open("postgres", psqlInfo)                                                                                       //bazaga ulanish
	if err != nil {
		fmt.Println("Error: Could not connect to the Postgres database")
	}
	err = db.Ping() //bazaga ulanishni tekshirish
	if err != nil {
		fmt.Println("Error: Could not establish a connection with the database")
		//panic(err)
	}
	fmt.Println("Successfully connected!") //bazaga ulanishni tekshirish
	return db                              //bazaga ulanishni qaytarish
}

func Login(c *gin.Context) {
	var user User
	_ = c.BindJSON(&user)
	db := connectDB()
	//if users table does not exist, create it and insert the user
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, access_token TEXT, id_token TEXT, ids TEXT, phone TEXT, email TEXT, password TEXT, name TEXT, photo_url TEXT, status TEXT, blocked BOOLEAN, role TEXT, region TEXT, device TEXT, created_at TIMESTAMP, updated_at TIMESTAMP)")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating the users table"})
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
		//_, err := db.Exec("INSERT INTO users (access_token, id_token, ids, phone, email, name, photo_url, status, blocked, role, region, device, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false, $9, $10, $11, $12, $13)", user.AccessToken, user.IdToken, user.Ids, user.Phone, user.Email, user.Name, user.PhotoUrl, user.Status, user.Role, user.Region, user.Device, user.CreatedAt, user.UpdatedAt) //insert the user
		_, err := db.Exec("INSERT INTO users (access_token, id_token, ids, phone, email, password, name, photo_url, blocked, role, region, device, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false, $9, $10, $11, $12, $13)", user.AccessToken, user.IdToken, user.Ids, user.Phone, user.Email, user.Password, user.Name, user.PhotoUrl, user.Role, user.Region, user.Device, user.CreatedAt, user.UpdatedAt) //insert the user
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating the user " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
		return
	} else {
		_, err := db.Exec("UPDATE users SET access_token=$1, id_token=$2, ids=$3, email=$4, name=$5, photo_url=$6, region=$7, device=$8, updated_at=$9 WHERE id=$10", user.AccessToken, user.IdToken, user.Ids, user.Email, user.Name, user.PhotoUrl, user.Region, user.Device, user.UpdatedAt, id) //update the user
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error updating the user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
		return
	}
}

func GetUsers(c *gin.Context) {
	db := connectDB()
	rows, err := db.Query("SELECT id, access_token, id_token, ids, phone, email, password, name, photo_url, blocked, role, region, device, created_at, updated_at FROM users")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
		return
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.AccessToken, &user.IdToken, &user.Ids, &user.Phone, &user.Email, &user.Password, &user.Name, &user.PhotoUrl, &user.Blocked, &user.Role, &user.Region, &user.Device, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the users"})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
	return
}
