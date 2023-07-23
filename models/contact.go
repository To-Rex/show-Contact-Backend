package models

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
