package services

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
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

type Contacts struct {
	DepartmentName        string `json:"department_name"`
	EmailAddresses        string `json:"email_addresses"`
	FamilyName            string `json:"family_name"`
	FullName              string `json:"full_name"`
	GivenName             string `json:"given_name"`
	JobTitle              string `json:"job_title"`
	MiddleName            string `json:"middle_name"`
	NamePrefix            string `json:"name_prefix"`
	NameSuffix            string `json:"name_suffix"`
	OrganizationName      string `json:"organization_name"`
	PhoneNumbers          string `json:"phone_numbers"`
	PhoneticFamilyName    string `json:"phonetic_family_name"`
	PhoneticGivenName     string `json:"phonetic_given_name"`
	PhoneticMiddleName    string `json:"phonetic_middle_name"`
	PhoneticOrganizationN string `json:"phonetic_organization_n"`
	PostalAddresses       string `json:"postal_addresses"`
	URLAddresses          string `json:"url_addresses"`
}
type ContactsList struct {
	Contacts []Contacts `json:"contacts"`
}

func ImPortContact(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	//save the file to the server
	err = c.SaveUploadedFile(file, "./files/EXcontacts.xlsx")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error saving the file"})
		return
	}

	f, err := excelize.OpenFile("./services/contacts.xlsx")
	if err != nil {
		c.JSON(400, gin.H{"error": "Error reading the file"})
		return
	}
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		c.JSON(400, gin.H{"error": "Error reading the file"})
		return
	}
	var contacts []Contacts
	for _, row := range rows {
		var contact Contacts
		for i, colCell := range row {
			switch rows[0][i] {
			case "Department Name":
				contact.DepartmentName = colCell
			case "Email Addresses":
				contact.EmailAddresses = colCell
			case "Family Name":
				contact.FamilyName = colCell
			case "Full Name":
				contact.FullName = colCell
			case "Given Name":
				contact.GivenName = colCell
			case "Job Title":
				contact.JobTitle = colCell
			case "Middle Name":
				contact.MiddleName = colCell
			case "Name Prefix":
				contact.NamePrefix = colCell
			case "Name Suffix":
				contact.NameSuffix = colCell
			case "Organization Name":
				contact.OrganizationName = colCell
			case "Phone Numbers":
				contact.PhoneNumbers = colCell
			case "Phonetic Family Name":
				contact.PhoneticFamilyName = colCell
			case "Phonetic Given Name":
				contact.PhoneticGivenName = colCell
			case "Phonetic Middle Name":
				contact.PhoneticMiddleName = colCell
			case "Phonetic Organization Name":
				contact.PhoneticOrganizationN = colCell
			case "Postal Addresses":
				contact.PostalAddresses = colCell
			case "URL Addresses":
				contact.URLAddresses = colCell
			}
		}
		contacts = append(contacts, contact)
	}
	c.JSON(200, gin.H{"contacts": contacts})
}

// ExportContact export contact
func ExportContact(c *gin.Context) {
	// Create a new spreadsheet and insert the default sheet.
	var contatcsList ContactsList
	db := connectDB()

	//if contacts table does not exist, create it and insert the user
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS contacts (department_name TEXT, email_addresses TEXT, family_name TEXT, full_name TEXT, given_name TEXT, job_title TEXT, middle_name TEXT, name_prefix TEXT, name_suffix TEXT, organization_name TEXT, phone_numbers TEXT, phonetic_family_name TEXT, phonetic_given_name TEXT, phonetic_middle_name TEXT, phonetic_organization_n TEXT, postal_addresses TEXT, url_addresses TEXT)")
	if err != nil {
		fmt.Println("Error: Could not create the contacts table")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating the contacts table"})
		return
	}
	//save db contactslist
	rows, err := db.Query("SELECT * FROM contacts")
	if err != nil {
		fmt.Println("Error: Could not get the contacts")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the contacts"})
		return
	}
	for rows.Next() {
		var contact Contacts
		err = rows.Scan(&contact.DepartmentName, &contact.EmailAddresses, &contact.FamilyName, &contact.FullName, &contact.GivenName, &contact.JobTitle, &contact.MiddleName, &contact.NamePrefix, &contact.NameSuffix, &contact.OrganizationName, &contact.PhoneNumbers, &contact.PhoneticFamilyName, &contact.PhoneticGivenName, &contact.PhoneticMiddleName, &contact.PhoneticOrganizationN, &contact.PostalAddresses, &contact.URLAddresses)
		if err != nil {
			fmt.Println("Error: Could not get the contacts")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting the contacts"})
			return
		}
		contatcsList.Contacts = append(contatcsList.Contacts, contact)
	}
	f := excelize.NewFile()
	// Create a new sheet.
	_, _ = f.NewSheet("Sheet1")
	// Set value of a cell.
	f.SetCellValue("Sheet1", "A1", "Department Name")
	f.SetCellValue("Sheet1", "B1", "Email Addresses")
	f.SetCellValue("Sheet1", "C1", "Family Name")
	f.SetCellValue("Sheet1", "D1", "Full Name")
	f.SetCellValue("Sheet1", "E1", "Given Name")
	f.SetCellValue("Sheet1", "F1", "Job Title")
	f.SetCellValue("Sheet1", "G1", "Middle Name")
	f.SetCellValue("Sheet1", "H1", "Name Prefix")
	f.SetCellValue("Sheet1", "I1", "Name Suffix")
	f.SetCellValue("Sheet1", "J1", "Organization Name")
	f.SetCellValue("Sheet1", "K1", "Phone Numbers")
	f.SetCellValue("Sheet1", "L1", "Phonetic Family Name")
	f.SetCellValue("Sheet1", "M1", "Phonetic Given Name")
	f.SetCellValue("Sheet1", "N1", "Phonetic Middle Name")
	f.SetCellValue("Sheet1", "O1", "Phonetic Organization Name")
	f.SetCellValue("Sheet1", "P1", "Postal Addresses")
	f.SetCellValue("Sheet1", "Q1", "URL Addresses")
	//contatcsList to excel file end
	//save the file to the server
	errs := f.SaveAs("./files/IMcontacts.xlsx")
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error saving the file"})
		return
	}
	contatcsList.Contacts = append(contatcsList.Contacts, Contacts{
		DepartmentName:        "Department Name",
		EmailAddresses:        "Email Addresses",
		FamilyName:            "Family Name",
		FullName:              "Full Name",
		GivenName:             "Given Name",
		JobTitle:              "Job Title",
		MiddleName:            "Middle Name",
		NamePrefix:            "Name Prefix",
		NameSuffix:            "Name Suffix",
		OrganizationName:      "Organization Name",
		PhoneNumbers:          "Phone Numbers",
		PhoneticFamilyName:    "Phonetic Family Name",
		PhoneticGivenName:     "Phonetic Given Name",
		PhoneticMiddleName:    "Phonetic Middle Name",
		PhoneticOrganizationN: "Phonetic Organization Name",
		PostalAddresses:       "Postal Addresses",
		URLAddresses:          "URL Addresses",
	})

	c.JSON(http.StatusOK, gin.H{"message": "Contact created successfully"})
}
