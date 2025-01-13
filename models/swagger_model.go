package models 

type CreateUser struct {
    Name  string `gorm:"not null" json:"Name"`
    Age   int    `json:"Age"`
    Email string `gorm:"unique" json:"Email"`
}