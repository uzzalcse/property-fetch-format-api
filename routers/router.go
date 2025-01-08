package routers

import (
    "property-fetch-format-api/controllers"
    beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Create namespace
	ns := beego.NewNamespace("/v1/api",
		beego.NSNamespace("/property",
			beego.NSRouter("/details/:propertyId", &controllers.PropertyDetailsController{}, "get:GetPropertyDetails"),
			beego.NSRouter("/:propertyId/gallery", &controllers.PropertyGalleryController{}, "get:GetPropertyGallery"),
		),
		beego.NSRouter("/propertyList", &controllers.PropertyListController{}, "get:GetPropertyList"),
	)

	// Register namespace
	beego.AddNamespace(ns)
}