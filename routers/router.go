// package routers

// import (
//     "property-fetch-format-api/controllers"
//     beego "github.com/beego/beego/v2/server/web"
// )

// func init() {
// 	// Create namespace
// 	ns := beego.NewNamespace("/v1/api",
// 		beego.NSNamespace("/property",
// 			beego.NSRouter("/details/:propertyId", &controllers.PropertyDetailsController{}, "get:GetPropertyDetails"),
// 			beego.NSRouter("/:propertyId/gallery", &controllers.PropertyGalleryController{}, "get:GetPropertyGallery"),
// 		),
// 		beego.NSRouter("/propertyList", &controllers.PropertyListController{}, "get:GetPropertyList"),
// 		beego.NSNamespace("/user",
// 		beego.NSRouter("/", &controllers.CreateUserController{}, "post:CreateUser"),
// 		beego.NSRouter("/:identifier", &controllers.UserController{}, "get:GetUser;put:UpdateUser;delete:DeleteUser"),
// 	),
		
// 	)

// 	// Register namespace
// 	beego.AddNamespace(ns)
// }



// package routers

// import (
//     "github.com/beego/beego/v2/server/web/swagger"
// 	 _"property-fetch-format-api/docs"
//     beego "github.com/beego/beego/v2/server/web"
//     "property-fetch-format-api/controllers"
// )

// func init() {
//     // Swagger Initialization
//     swagger.SwaggerInfo.Title = "Property Fetch Format API"
//     swagger.SwaggerInfo.Description = "API for managing properties and users"
//     swagger.SwaggerInfo.Version = "1.0"
//     swagger.SwaggerInfo.Host = "localhost:8080"
//     swagger.SwaggerInfo.BasePath = "/v1/api"
//     swagger.SwaggerInfo.Schemes = []string{"http", "https"}

//     // Create namespace
//     ns := beego.NewNamespace("/v1/api",
//         beego.NSNamespace("/property",
//             beego.NSRouter("/details/:propertyId", &controllers.PropertyDetailsController{}, "get:GetPropertyDetails"),
//             beego.NSRouter("/:propertyId/gallery", &controllers.PropertyGalleryController{}, "get:GetPropertyGallery"),
//         ),
//         beego.NSRouter("/propertyList", &controllers.PropertyListController{}, "get:GetPropertyList"),
//         beego.NSNamespace("/user",
//             beego.NSRouter("/", &controllers.CreateUserController{}, "post:CreateUser"),
//             beego.NSRouter("/:identifier", &controllers.UserController{}, "get:GetUser;put:UpdateUser;delete:DeleteUser"),
//         ),
//     )

//     // Register namespace
//     beego.AddNamespace(ns)

//     // Enable Swagger endpoint
//     beego.Router("/swagger/*", swagger.WrapHandler)
// }


package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/swaggo/http-swagger"  // Import the http-swagger package
	_ "property-fetch-format-api/docs"  // Import generated Swagger docs
	"property-fetch-format-api/controllers"
)

// SwaggerController will serve the Swagger UI using the standard HTTP handler
type SwaggerController struct {
	web.Controller
}

func (c *SwaggerController) Get() {
	httpSwagger.Handler().ServeHTTP(c.Ctx.ResponseWriter, c.Ctx.Request)
}

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
	web.Router("/swagger/*", &SwaggerController{})
}
