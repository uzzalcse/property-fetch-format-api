package services

import (
    "fmt"
    "regexp"
    "strings"

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
        return fmt.Errorf("validation error: name is required")
    }

    if user.Age <= 0 {
        return fmt.Errorf("validation error: insert a positive integer for age field")
    }

    if !s.isValidEmail(user.Email) {
        return fmt.Errorf("validation error: invalid email format")
    }

    return nil
}

// isValidEmail checks if the given email has a valid format
func (s *UserService) isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return re.MatchString(strings.ToLower(email))
}

// GetUserByIdentifier retrieves a user by ID or email
func (s *UserService) GetUserByIdentifier(identifier string) (*models.User, error) {
    var user models.User
    db := dao.GetDB()

    if s.isValidEmail(identifier) {
        if err := db.Where("email = ?", identifier).First(&user).Error; err != nil {
            return nil, fmt.Errorf("user not found")
        }
    } else {
        if err := db.Where("id = ?", identifier).First(&user).Error; err != nil {
            return nil, fmt.Errorf("user not found")
        }
    }

    return &user, nil
}

// UpdateUserByIdentifier updates a user by ID or email
func (s *UserService) UpdateUserByIdentifier(identifier string, userUpdate *models.User) (*models.User, error) {
    var user models.User
    db := dao.GetDB()

    if s.isValidEmail(identifier) {
        if err := db.Where("email = ?", identifier).First(&user).Error; err != nil {
            return nil, fmt.Errorf("user not found")
        }
    } else {
        if err := db.Where("id = ?", identifier).First(&user).Error; err != nil {
            return nil, fmt.Errorf("user not found")
        }
    }

    if userUpdate.Name != "" {
        user.Name = userUpdate.Name
    }
    if userUpdate.Age > 0 {
        user.Age = userUpdate.Age
    }
    if s.isValidEmail(userUpdate.Email) {
        user.Email = userUpdate.Email
    }

    if err := db.Save(&user).Error; err != nil {
        return nil, fmt.Errorf("failed to update user")
    }

    return &user, nil
}

// DeleteUserByIdentifier deletes a user by ID or email
func (s *UserService) DeleteUserByIdentifier(identifier string) error {
    var user models.User
    db := dao.GetDB()

    if s.isValidEmail(identifier) {
        if err := db.Where("email = ?", identifier).First(&user).Error; err != nil {
            return fmt.Errorf("user not found")
        }
    } else {
        if err := db.Where("id = ?", identifier).First(&user).Error; err != nil {
            return fmt.Errorf("user not found")
        }
    }

    if err := db.Delete(&user).Error; err != nil {
        return fmt.Errorf("failed to delete user")
    }

    return nil
}