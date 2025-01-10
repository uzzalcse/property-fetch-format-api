package main

import (
	_ "property-fetch-format-api/routers"
	"property-fetch-format-api/db_services"
	"log"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	db, err := db_services.InitDB()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
	beego.Run()
}

