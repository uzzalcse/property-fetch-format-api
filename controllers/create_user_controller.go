// create_user_controller.go
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

func (u *CreateUserController) CreateUser() {
    var user models.User

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

    if err := json.Unmarshal(body, &user); err != nil {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]interface{}{
            "status": http.StatusBadRequest,
            "error": fmt.Sprintf("Failed to unmarshal JSON: %v", err),
        }
        u.ServeJSON()
        return
    }

    errChan := make(chan error)
    
    go func() {
        userService := services.UserService{}
        errChan <- userService.CreateUser(&user)
    }()

    if err := <-errChan; err != nil {
        statusCode := http.StatusInternalServerError
        if err.Error() == "email already exists" {
            statusCode = http.StatusConflict
        }
        u.Ctx.Output.SetStatus(statusCode)
        u.Data["json"] = map[string]interface{}{
            "status": statusCode,
            "error": err.Error(),
        }
        u.ServeJSON()
        return
    }

    u.Ctx.Output.SetStatus(http.StatusCreated)
    u.Data["json"] = map[string]interface{}{
        "status": http.StatusCreated,
        "message": "User created successfully",
        "data": user,
    }
    u.ServeJSON()
}