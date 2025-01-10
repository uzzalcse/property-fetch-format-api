package controllers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"

    "property-fetch-format-api/models"
    "property-fetch-format-api/services"

    "github.com/beego/beego/v2/server/web"
)

type CreateUserController struct {
    web.Controller
}

// CreateUser handles the creation of a new user
func (u *CreateUserController) CreateUser() {
    var user models.User

    // Read the request body
    body, err := io.ReadAll(u.Ctx.Request.Body)
    if err != nil {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to read request body: %v", err)}
        u.ServeJSON()
        return
    }

    // Bind the JSON payload to the user struct
    if err := json.Unmarshal(body, &user); err != nil {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to unmarshal JSON: %v", err)}
        u.ServeJSON()
        return
    }

    // Use the service layer to create the user
    userService := services.UserService{}
    if err := userService.CreateUser(&user); err != nil {
        u.Ctx.Output.SetStatus(http.StatusInternalServerError)
        u.Data["json"] = map[string]string{"error": err.Error()}
        u.ServeJSON()
        return
    }

    // Return the created user as JSON response
    u.Data["json"] = user
    u.ServeJSON()
}