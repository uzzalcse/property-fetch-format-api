package db_services

import (
    "github.com/jinzhu/gorm"
	"property-fetch-format-api/models"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

func InitDB() (*gorm.DB, error) {
    // Use the Docker service name 'db' as the host
    db, err := gorm.Open("postgres", "host=db user=postgres dbname=property_user_db sslmode=disable password=postgres")
    if err != nil {
        return nil, err
    }
	 // Drop the existing table
	//db.DropTableIfExists(&models.User{})
    db.AutoMigrate(&models.User{})
    return db, nil
}