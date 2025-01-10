package models

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
    ID    uint   `gorm:"primary_key" json:"id"`
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}

// TableName sets the table name for the User struct
func (User) TableName() string {
    return "property_users"
}

func InitDB() (*gorm.DB, error) {
    // Use the Docker service name 'db' as the host
    db, err := gorm.Open("postgres", "host=db user=postgres dbname=property_user_db sslmode=disable password=postgres")
    if err != nil {
        return nil, err
    }
    db.AutoMigrate(&User{})
    return db, nil
}