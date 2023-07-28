package config

import (
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	host     = "containers-us-west-73.railway.app" //host
	port     = 6725                                //port
	user     = "postgres"                          //foydalanuvchi
	password = "RHnJjy7DZt2HGGw7kn28"              //parol
	dbname   = "railway"                           //baza nomi
)

func ConnectDB() *sql.DB { //Dastur bilan bazaga ulanish
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname) //bazaga ulanish uchun ma'lumotlar
	db, err := sql.Open("postgres", psqlInfo)                                                                                       //bazaga ulanish
	if err != nil {
		//fmt.Println("Error: Could not connect to the Postgres database")
		return nil
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error: Could not establish a connection with the database")
		//panic(err)
		return nil
	}
	return db
}

func checkPasswordHash(password, hash string) bool { //parolni tekshirish
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) //parolni tekshirish
	return err == nil                                                    //agar xato bo'lmasa true qaytaradi
}

func PasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash)
}

func GenerateToken(email string, name string, roles string) (string, error) { //token yaratish
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //token yaratish
		"email": email,
		"name":  name,
		//tugash vaqti 24 soat keyin
		"created_at": time.Now().Add(time.Hour * 24).Unix(),
		"roles":      roles,
	})
	tokenString, err := token.SignedString([]byte("secret")) //tokenni shifrlash
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
