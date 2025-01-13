package routers

import (
	"github.com/beego/beego/v2/server/web"
	_ "property-fetch-format-api/docs"  // Import generated Swagger docs
	"property-fetch-format-api/controllers"
)

func init() {
	// Define API namespaces and routes
	ns := web.NewNamespace("/v1/api",
		web.NSNamespace("/property",
			web.NSRouter("/details/:propertyId", &controllers.PropertyDetailsController{}, "get:GetPropertyDetails"),
			web.NSRouter("/:propertyId/gallery", &controllers.PropertyGalleryController{}, "get:GetPropertyGallery"),
		),
		web.NSRouter("/propertyList", &controllers.PropertyListController{}, "get:GetPropertyList"),
		web.NSNamespace("/user",
			web.NSRouter("/", &controllers.CreateUserController{}, "post:CreateUser"),
			web.NSRouter("/:identifier", &controllers.UserController{}, "get:GetUser;put:UpdateUser;delete:DeleteUser"),
		),
	)

	// Register namespace
	web.AddNamespace(ns)

	// Register the Swagger UI endpoint
	web.Router("/swagger/*", &controllers.SwaggerController{})
}