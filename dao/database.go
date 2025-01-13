package dao

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "property-fetch-format-api/models"
)

var db *gorm.DB

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
    var err error
    db, err = gorm.Open("postgres", "host=db user=postgres dbname=property_user_db sslmode=disable password=postgres")
    if err != nil {
        return nil, err
    }

    // Auto migrate the User model
    db.AutoMigrate(&models.User{})

    return db, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
    return db
}