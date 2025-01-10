package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"property-fetch-format-api/dao"
	"property-fetch-format-api/models"

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

    // Log the received JSON payload
    fmt.Println("Received JSON:", string(body))

    // Bind the JSON payload to the user struct
    if err := json.Unmarshal(body, &user); err != nil {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to unmarshal JSON: %v", err)}
        u.ServeJSON()
        return
    }

    // Validate the user data
    if err := validateUser(&user); err != nil {
        u.Ctx.Output.SetStatus(http.StatusBadRequest)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("Validation error: %v", err)}
        u.ServeJSON()
        return
    }

    // Create the new user in the database
    if err := dao.GetDB().Create(&user).Error; err != nil {
        u.Ctx.Output.SetStatus(http.StatusInternalServerError)
        u.Data["json"] = map[string]string{"error": fmt.Sprintf("Database error: %v", err)}
        u.ServeJSON()
        return
    }

    // Return the created user as JSON response
    u.Data["json"] = user
    u.ServeJSON()
}

// validateUser performs basic validation on the user struct
func validateUser(user *models.User) error {
    if user.Name == "" {
        return fmt.Errorf("name is required")
    }

    if user.Age <= 0 {
        return fmt.Errorf("age must be a positive integer")
    }

    if !isValidEmail(user.Email) {
        return fmt.Errorf("invalid email format")
    }

    return nil
}

// isValidEmail checks if the given email has a valid format
func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return re.MatchString(email)
}