package main

import (
	_ "property-fetch-format-api/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

