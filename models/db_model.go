package models

type User struct {
    ID    uint   `gorm:"primary_key;auto_increment" json:"id"`
    Name  string `gorm:"not null" json:"name"`
    Age   int    `json:"age"`
    Email string `gorm:"unique" json:"email"`
}

// TableName sets the table name for the User struct
func (User) TableName() string {
    return "property_users"
}
