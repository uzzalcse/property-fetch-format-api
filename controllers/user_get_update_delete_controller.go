// user_get_update_delete_controller.go
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

// @Title GetUser
// @Description Retrieve user by ID or email
// @Param identifier path string true "User ID (integer) or email (string)"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]interface{} "User not found"
// @router /v1/api/user/{identifier} [get]
func (u *UserController) GetUser() {
    identifier := u.Ctx.Input.Param(":identifier")
    
    resultChan := make(chan struct {
        user *models.User
        err  error
    })
    
    go func() {
        userService := services.UserService{}
        user, err := userService.GetUserByIdentifier(identifier)
        resultChan <- struct {
            user *models.User
            err  error
        }{user, err}
    }()
    
    result := <-resultChan
    if result.err != nil {
        u.Ctx.Output.SetStatus(http.StatusNotFound)
        u.Data["json"] = map[string]interface{}{
            "status": http.StatusNotFound,
            "error": fmt.Sprintf("User not found: %v", result.err),
        }
        u.ServeJSON()
        return
    }

    u.Ctx.Output.SetStatus(http.StatusOK)
    u.Data["json"] = map[string]interface{}{
        "status": http.StatusOK,
        "message": "User retrieved successfully",
        "data": result.user,
    }
    u.ServeJSON()
}


// @Title UpdateUser
// @Description Update user details
// @Param identifier path string true "User ID or email"
// @Param body body models.User true "Updated user details"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{} "Validation error"
// @router /v1/api/user/{identifier} [put]
func (u *UserController) UpdateUser() {
    identifier := u.Ctx.Input.Param(":identifier")
    var userUpdate models.User

    body, err := io.ReadAll(u.Ctx.Request.Body)
    if err != nil {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]interface{}{
            "status": http.StatusBadRequest,
            "error": fmt.Sprintf("Failed to read request body: %v", err),
        }
        u.ServeJSON()
        return
    }

    if err := json.Unmarshal(body, &userUpdate); err != nil {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]interface{}{
            "status": http.StatusBadRequest,
            "error": fmt.Sprintf("Failed to unmarshal JSON: %v", err),
        }
        u.ServeJSON()
        return
    }

    if userUpdate.Name == "" && userUpdate.Age == 0 && userUpdate.Email == "" {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]interface{}{
            "status": http.StatusBadRequest,
            "error": "At least one field must be provided for update",
        }
        u.ServeJSON()
        return
    }

    resultChan := make(chan struct {
        user *models.User
        err  error
    })

    go func() {
        userService := services.UserService{}
        user, err := userService.UpdateUserByIdentifier(identifier, &userUpdate)
        resultChan <- struct {
            user *models.User
            err  error
        }{user, err}
    }()

    result := <-resultChan
    if result.err != nil {
        u.Ctx.Output.SetStatus(http.StatusInternalServerError)
        u.Data["json"] = map[string]interface{}{
            "status": http.StatusInternalServerError,
            "error": fmt.Sprintf("Failed to update user: %v", result.err),
        }
        u.ServeJSON()
        return
    }

    u.Ctx.Output.SetStatus(http.StatusOK)
    u.Data["json"] = map[string]interface{}{
        "status": http.StatusOK,
        "message": "User updated successfully",
        "data": result.user,
    }
    u.ServeJSON()
}


// @Title DeleteUser
// @Description Delete user by ID or email
// @Param identifier path string true "User ID or email"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]interface{} "User not found"
// @router /v1/api/user/{identifier} [delete]
func (u *UserController) DeleteUser() {
    identifier := u.Ctx.Input.Param(":identifier")

    // Get user before deletion
    userChan := make(chan struct {
        user *models.User
        err  error
    })
    
    go func() {
        userService := services.UserService{}
        user, err := userService.GetUserByIdentifier(identifier)
        userChan <- struct {
            user *models.User
            err  error
        }{user, err}
    }()

    userResult := <-userChan
    if userResult.err != nil {
        u.Ctx.Output.SetStatus(http.StatusNotFound)
        u.Data["json"] = map[string]interface{}{
            "status": http.StatusNotFound,
            "error": fmt.Sprintf("User not found: %v", userResult.err),
        }
        u.ServeJSON()
        return
    }

    errChan := make(chan error)
    go func() {
        userService := services.UserService{}
        errChan <- userService.DeleteUserByIdentifier(identifier)
    }()

    if err := <-errChan; err != nil {
        u.Ctx.Output.SetStatus(http.StatusInternalServerError)
        u.Data["json"] = map[string]interface{}{
            "status": http.StatusInternalServerError,
            "error": fmt.Sprintf("Failed to delete user: %v", err),
        }
        u.ServeJSON()
        return
    }

    u.Ctx.Output.SetStatus(http.StatusOK)
    u.Data["json"] = map[string]interface{}{
        "status": http.StatusOK,
        "message": "User deleted successfully",
        "data": userResult.user,
    }
    u.ServeJSON()
}