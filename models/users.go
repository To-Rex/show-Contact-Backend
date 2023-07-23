package models

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
	Token       string `json:"token"`
}
