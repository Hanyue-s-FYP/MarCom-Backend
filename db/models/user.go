package models

type User struct {
	ID          int
	Username    string
	Password    string
	DisplayName string
	Email       string
	Status      int
	PhoneNumber string
}

type Investor struct {
	User
	ID int
}

type Business struct {
	ID           int
	Description  string
	BusinessType string
	CoverImgPath string
}
