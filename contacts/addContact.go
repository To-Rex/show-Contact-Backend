package contacts

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Contact struct {
	DisplayName           string `json:"display_name"`
	GivenName             string `json:"given_name"`
	MiddleName            string `json:"middle_name"`
	Prefix                string `json:"prefix"`
	Suffix                string `json:"suffix"`
	FamilyName            string `json:"family_name"`
	Company               string `json:"company"`
	JobTitle              string `json:"job_title"`
	Emails                string `json:"emails"`
	Phones                string `json:"phones"`
	PostalAddresses       string `json:"postal_addresses"`
	Avatar                byte   `json:"avatar"`
	Birthday              string `json:"birthday"`
	AndroidAccountType    string `json:"android_account_type"`
	AndroidAccountTypeRaw string `json:"android_account_type_raw"`
	AndroidAccountName    string `json:"android_account_name"`
}

type Contacts struct {
	Contacts []Contact `json:"contacts"`
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

// add contact to database
func AddContact(c *gin.Context) {
	/*{
			"contact":
	[
		{
		"display_name":"Dilshodjon",
		"given_name":"",
		"prefix":"",
		"suffix":"",
		"family_name":"",
		"company":"Staj",
		"job_title":"",
		"emails":[],
		"phones":["+998995340313"],
		"postal_addresses":[],
		"avatar":[],
		"birthday":"",
		"android_account_type":"AndroidAccountType.google",
		"android_account_type_raw":"com.google",
		"android_account_name":"dev.dilshodjon@gmail.com"
		},
		{
		"display_name":"Dilshodjon1",
		"given_name":"",
		"prefix":"",
		"suffix":"",
		"family_name":"",
		"company":"Staj",
		"job_title":"",
		"emails":[],
		"phones":["+998995340314"],
		"postal_addresses":[],
		"avatar":[],
		"birthday":"",
		"android_account_type":"AndroidAccountType.google",
		"android_account_type_raw":"com.google",
		"android_account_name":"dev.dilshodjon@gmail.com"
		},
		{
		"display_name":"Dilshodjon2",
		"given_name":"",
		"prefix":"",
		"suffix":"",
		"family_name":"",
		"company":"Staj",
		"job_title":"",
		"emails":[],
		"phones":["+998995340315"],
		"postal_addresses":[],
		"avatar":[],
		"birthday":"",
		"android_account_type":"AndroidAccountType.google",
		"android_account_type_raw":"com.google",
		"android_account_name":"dev.dilshodjon@gmail.com"
		}
		]
		}*/
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	db := connectDB()
	var contacts Contacts
	_ = c.BindJSON(&contacts)
	if len(contacts.Contacts) == 0 {
		c.JSON(400, gin.H{"error": "contacts is empty"})
		return
	}
	//fmt.Println(contacts.length)
	//if database is empty then create table contacts and insert contacts
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS contacts (id SERIAL PRIMARY KEY, display_name TEXT, given_name TEXT, middle_name TEXT, prefix TEXT, suffix TEXT, family_name TEXT, company TEXT, job_title TEXT, emails TEXT, phones TEXT, postal_addresses TEXT, avatar BYTEA, birthday TEXT, android_account_type TEXT, android_account_type_raw TEXT, android_account_name TEXT)")
	if err != nil {
		c.JSON(400, gin.H{"error": "Error creating the contacts table"})
		return
	}
	for _, contact := range contacts.Contacts {
		if contact.Phones != "" && contact.DisplayName != "" {
			err = db.QueryRow("SELECT id FROM contacts WHERE phones=$1 AND display_name=$2", contact.Phones, contact.DisplayName).Scan(&contact.DisplayName)
			fmt.Println(err)
			if err != nil {
				_, err = db.Exec("INSERT INTO contacts (display_name, given_name, middle_name, prefix, suffix, family_name, company, job_title, emails, phones, postal_addresses, avatar, birthday, android_account_type, android_account_type_raw, android_account_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ,$11, $12, $13, $14, $15, $16)", contact.DisplayName, contact.GivenName, contact.MiddleName, contact.Prefix, contact.Suffix, contact.FamilyName, contact.Company, contact.JobTitle, contact.Emails, contact.Phones, contact.PostalAddresses, contact.Avatar, contact.Birthday, contact.AndroidAccountType, contact.AndroidAccountTypeRaw, contact.AndroidAccountName)
				if err != nil {
					c.JSON(400, gin.H{"error": "Error inserting the contact"})
					return
				}
			}
		}
	}
	defer db.Close()
	c.JSON(200, gin.H{"message": "Successfully added contacts"})
}
