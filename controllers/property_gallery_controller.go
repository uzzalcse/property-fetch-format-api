package controllers

import (
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyGalleryController struct {
	beego.Controller
}

func (c *PropertyGalleryController) GetPropertyGallery() {
	c.Data["json"] = map[string]string{"propertyId": c.Ctx.Input.Param(":propertyId")}
	c.ServeJSON()
}