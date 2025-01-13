package services

import (
    "fmt"
    "regexp"
    "strings"
    "sync"

    "property-fetch-format-api/dao"
    "property-fetch-format-api/models"
)

type UserService struct{}

func (s *UserService) CreateUser(user *models.User) error {
    // Check if email already exists
    var existingUser models.User
    if err := dao.GetDB().Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
        return fmt.Errorf("email already exists")
    }

    errChan := make(chan error)
    
    go func() {
        if err := s.validateUser(user); err != nil {
            errChan <- err
            return
        }

        errChan <- dao.GetDB().Create(user).Error
    }()

    return <-errChan
}

func (s *UserService) validateUser(user *models.User) error {
    var wg sync.WaitGroup
    errChan := make(chan error, 3)

    wg.Add(3)
    
    go func() {
        defer wg.Done()
        if user.Name == "" {
            errChan <- fmt.Errorf("validation error: name is required")
        } else if len(user.Name) < 2 {
            errChan <- fmt.Errorf("validation error: name must be at least 2 characters long")
        }
    }()

    go func() {
        defer wg.Done()
        if user.Age <= 0 {
            errChan <- fmt.Errorf("validation error: age must be a positive integer")
        } else if user.Age > 150 {
            errChan <- fmt.Errorf("validation error: age must be less than 150")
        }
    }()

    go func() {
        defer wg.Done()
        if user.Email == "" {
            errChan <- fmt.Errorf("validation error: email is required")
        } else if !s.isValidEmail(user.Email) {
            errChan <- fmt.Errorf("validation error: invalid email format")
        }
    }()

    go func() {
        wg.Wait()
        close(errChan)
    }()

    for err := range errChan {
        if err != nil {
            return err
        }
    }

    return nil
}

func (s *UserService) isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return re.MatchString(strings.ToLower(email))
}

func (s *UserService) GetUserByIdentifier(identifier string) (*models.User, error) {
    var user models.User
    resultChan := make(chan struct {
        user models.User
        err  error
    })

    go func() {
        db := dao.GetDB()
        var err error

        if s.isValidEmail(identifier) {
            err = db.Where("email = ?", identifier).First(&user).Error
        } else {
            err = db.Where("id = ?", identifier).First(&user).Error
        }

        if err != nil {
            err = fmt.Errorf("user not found")
        }

        resultChan <- struct {
            user models.User
            err  error
        }{user, err}
    }()

    result := <-resultChan
    if result.err != nil {
        return nil, result.err
    }
    return &result.user, nil
}

func (s *UserService) UpdateUserByIdentifier(identifier string, userUpdate *models.User) (*models.User, error) {
    var user models.User
    resultChan := make(chan struct {
        user *models.User
        err  error
    })

    go func() {
        db := dao.GetDB()
        var err error

        if s.isValidEmail(identifier) {
            err = db.Where("email = ?", identifier).First(&user).Error
        } else {
            err = db.Where("id = ?", identifier).First(&user).Error
        }

        if err != nil {
            resultChan <- struct {
                user *models.User
                err  error
            }{nil, fmt.Errorf("user not found")}
            return
        }

        // Validate email if it's being updated
        if userUpdate.Email != "" && userUpdate.Email != user.Email {
            var existingUser models.User
            if err := db.Where("email = ?", userUpdate.Email).First(&existingUser).Error; err == nil {
                resultChan <- struct {
                    user *models.User
                    err  error
                }{nil, fmt.Errorf("email already exists")}
                return
            }
        }

        if userUpdate.Name != "" {
            if len(userUpdate.Name) < 2 {
                resultChan <- struct {
                    user *models.User
                    err  error
                }{nil, fmt.Errorf("validation error: name must be at least 2 characters long")}
                return
            }
            user.Name = userUpdate.Name
        }
        if userUpdate.Age > 0 {
            if userUpdate.Age > 150 {
                resultChan <- struct {
                    user *models.User
                    err  error
                }{nil, fmt.Errorf("validation error: age must be less than 150")}
                return
            }
            user.Age = userUpdate.Age
        }
        if userUpdate.Email != "" && s.isValidEmail(userUpdate.Email) {
            user.Email = userUpdate.Email
        }

        if err := db.Save(&user).Error; err != nil {
            resultChan <- struct {
                user *models.User
                err  error
            }{nil, fmt.Errorf("failed to update user")}
            return
        }

        resultChan <- struct {
            user *models.User
            err  error
        }{&user, nil}
    }()

    result := <-resultChan
    return result.user, result.err
}

func (s *UserService) DeleteUserByIdentifier(identifier string) error {
    errChan := make(chan error)

    go func() {
        var user models.User
        db := dao.GetDB()
        var err error

        if s.isValidEmail(identifier) {
            err = db.Where("email = ?", identifier).First(&user).Error
        } else {
            err = db.Where("id = ?", identifier).First(&user).Error
        }

        if err != nil {
            errChan <- fmt.Errorf("user not found")
            return
        }

        errChan <- db.Delete(&user).Error
    }()

    return <-errChan
}