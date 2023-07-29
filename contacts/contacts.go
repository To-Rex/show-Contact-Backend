package contacts

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"show_contact/config"
	"show_contact/models"
	"strings"
	"time"
)

func SearchContact(c *gin.Context) {
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
	searchTerm := c.Query("search")
	println(searchTerm)
	contacts, err := searchContacts(db, searchTerm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the contacts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}

func searchContacts(db *sql.DB, searchTerm string) ([]models.Contact, error) {
	rows, err := db.Query("SELECT display_name, given_name, middle_name, prefix, suffix, family_name, company, job_title, emails, phones, postal_addresses, avatar, birthday, android_account_type, android_account_type_raw, android_account_name FROM contacts WHERE display_name LIKE $1 OR given_name LIKE $1 OR middle_name LIKE $1 OR prefix LIKE $1 OR suffix LIKE $1 OR family_name LIKE $1 OR company LIKE $1 OR job_title LIKE $1 OR emails LIKE $1 OR phones LIKE $1 OR postal_addresses LIKE $1 OR android_account_type LIKE $1 OR android_account_type_raw LIKE $1 OR android_account_name LIKE $1", "%"+searchTerm+"%")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	var contacts []models.Contact
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(&contact.DisplayName, &contact.GivenName, &contact.MiddleName, &contact.Prefix, &contact.Suffix, &contact.FamilyName, &contact.Company, &contact.JobTitle, &contact.Emails, &contact.Phones, &contact.PostalAddresses, &contact.Avatar, &contact.Birthday, &contact.AndroidAccountType, &contact.AndroidAccountTypeRaw, &contact.AndroidAccountName)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}
