package services

import (
    "fmt"
    "regexp"

    "property-fetch-format-api/dao"
    "property-fetch-format-api/models"
)

// UserService handles the business logic for user operations
type UserService struct{}

// CreateUser creates a new user in the database
func (s *UserService) CreateUser(user *models.User) error {
    // Validate the user data
    if err := s.validateUser(user); err != nil {
        return err
    }

    // Create the new user in the database
    if err := dao.GetDB().Create(user).Error; err != nil {
        return fmt.Errorf("database error: %v", err)
    }

    return nil
}

// validateUser performs basic validation on the user struct
func (s *UserService) validateUser(user *models.User) error {
    if user.Name == "" {
        return fmt.Errorf("name is required")
    }

    if user.Age <= 0 {
        return fmt.Errorf("age must be a positive integer")
    }

    if !s.isValidEmail(user.Email) {
        return fmt.Errorf("invalid email format")
    }

    return nil
}

// isValidEmail checks if the given email has a valid format
func (s *UserService) isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return re.MatchString(email)
}