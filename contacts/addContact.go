package contacts

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"show_contact/config"
	"show_contact/models"
)

func AddContact(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	db := config.ConnectDB()
	var contacts models.Contacts
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
