package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type PropertyListController struct {
	beego.Controller
}

func (c *PropertyListController) GetPropertyList() {
	// Get property IDs from query parameter
	propertyIDs := c.GetString("propertyIds")

	// Split the IDs by comma
	ids := strings.Split(propertyIDs, ",")

	// Create response data
	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"propertyIds": ids,
		},
	}

	// Send JSON response
	c.Data["json"] = response
	c.ServeJSON()
}