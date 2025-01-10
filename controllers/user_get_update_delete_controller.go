package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"property-fetch-format-api/models"
	"property-fetch-format-api/services"

	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
    beego.Controller
}

// GetUser handles getting a user by ID or email
func (u *UserController) GetUser() {
    identifier := u.Ctx.Input.Param(":identifier")

    userService := services.UserService{}
    user, err := userService.GetUserByIdentifier(identifier)
    if err != nil {
        u.Ctx.Output.SetStatus(http.StatusNotFound)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("User not found: %v", err)}
        u.ServeJSON()
        return
    }

    u.Data["json"] = user
    u.ServeJSON()
}

// UpdateUser handles updating a user by ID or email
func (u *UserController) UpdateUser() {
    identifier := u.Ctx.Input.Param(":identifier")
    var userUpdate models.User

    // Read the request body
    body, err := io.ReadAll(u.Ctx.Request.Body)
    if err != nil {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to read request body: %v", err)}
        u.ServeJSON()
        return
    }

    // Log the received JSON payload
    fmt.Println("Received JSON:", string(body))

    // Bind the JSON payload to the user struct
    if err := json.Unmarshal(body, &userUpdate); err != nil {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to unmarshal JSON: %v", err)}
        u.ServeJSON()
        return
    }

    userService := services.UserService{}
    user, err := userService.UpdateUserByIdentifier(identifier, &userUpdate)
    if err != nil {
        u.Ctx.Output.SetStatus(http.StatusInternalServerError)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to update user: %v", err)}
        u.ServeJSON()
        return
    }

    u.Data["json"] = user
    u.ServeJSON()
}

// DeleteUser handles deleting a user by ID or email
func (u *UserController) DeleteUser() {
    identifier := u.Ctx.Input.Param(":identifier")

    userService := services.UserService{}
    if err := userService.DeleteUserByIdentifier(identifier); err != nil {
        u.Ctx.Output.SetStatus(http.StatusInternalServerError)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to delete user: %v", err)}
        u.ServeJSON()
        return
    }

    u.Ctx.Output.SetStatus(http.StatusNoContent)
    u.ServeJSON()
}