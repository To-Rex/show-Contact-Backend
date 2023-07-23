package config

import (
	"database/sql"
	"fmt"
)

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
