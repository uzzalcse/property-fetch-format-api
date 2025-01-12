package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"property-fetch-format-api/controllers"
	"property-fetch-format-api/models"
	"property-fetch-format-api/dao"
	"testing"

	"github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (suite *UserControllerTestSuite) SetupSuite() {
	dao.InitDB()
	
	ns := web.NewNamespace("/v1/api",
		web.NSNamespace("/user",
			web.NSRouter("/", &controllers.CreateUserController{}, "post:CreateUser"),
			web.NSRouter("/:identifier", &controllers.UserController{}, "get:GetUser;put:UpdateUser;delete:DeleteUser"),
		),
	)
	web.AddNamespace(ns)
}

func (suite *UserControllerTestSuite) SetupTest() {
	db := dao.GetDB()
	db.Exec("DELETE FROM property_users")
}

func (suite *UserControllerTestSuite) TearDownTest() {
	db := dao.GetDB()
	db.Exec("DELETE FROM property_users")
}

func (suite *UserControllerTestSuite) TestCreateUserController() {
	tests := []struct {
		name           string
		input         models.User
		setupFunc     func()
		expectedCode  int
		expectedError string
	}{
		{
			name: "Valid User Creation",
			input: models.User{
				Name:  "John Doe",
				Email: "john.khan@example.com",
				Age:   30,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Invalid Email Format",
			input: models.User{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   30,
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "validation error: invalid email format",
		},
		{
			name: "Empty Name",
			input: models.User{
				Email: "empty.name@example.com",
				Age:   30,
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "validation error: name is required",
		},
		{
			name: "Invalid Age",
			input: models.User{
				Name:  "John Doe",
				Email: "invalid.age@example.com",
				Age:   -1,
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "validation error: age must be a positive integer",
		},
		{
			name: "Duplicate Email",
			input: models.User{
				Name:  "John Doe",
				Email: "existing@example.com",
				Age:   30,
			},
			setupFunc: func() {
				user := models.User{
					Name:  "Existing User",
					Email: "existing@example.com",
					Age:   25,
				}
				result := dao.GetDB().Create(&user)
				if result.Error != nil {
					panic(result.Error)
				}
			},
			expectedCode:  http.StatusConflict,
			expectedError: "email already exists",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			jsonBody, err := json.Marshal(tt.input)
			assert.NoError(suite.T(), err)

			req := httptest.NewRequest("POST", "/v1/api/user", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			web.BeeApp.Handlers.ServeHTTP(w, req)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), tt.expectedCode, w.Code)

			if tt.expectedError != "" {
				assert.Equal(suite.T(), tt.expectedError, response["error"])
			} else if w.Code == http.StatusCreated {
				assert.Equal(suite.T(), "User created successfully", response["message"])
				if userData, ok := response["data"].(map[string]interface{}); ok {
					assert.Equal(suite.T(), tt.input.Name, userData["name"])
					assert.Equal(suite.T(), tt.input.Email, userData["email"])
					assert.Equal(suite.T(), float64(tt.input.Age), userData["age"])
				}
			}
		})
	}
}

func (suite *UserControllerTestSuite) TestGetUserController() {
	testUser := models.User{
		Name:  "Test User",
		Email: "test.get@example.com",
		Age:   30,
	}
	result := dao.GetDB().Create(&testUser)
	assert.NoError(suite.T(), result.Error)

	tests := []struct {
		name          string
		identifier    string
		expectedCode  int
		expectedError string
	}{
		{
			name:         "Get User By Email",
			identifier:   "test.get@example.com",
			expectedCode: http.StatusOK,
		},
		{
			name:          "User Not Found",
			identifier:    "nonexistent@example.com",
			expectedCode:  http.StatusNotFound,
			expectedError: "User not found: user not found",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/v1/api/user/"+tt.identifier, nil)
			w := httptest.NewRecorder()
			web.BeeApp.Handlers.ServeHTTP(w, req)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Equal(suite.T(), tt.expectedError, response["error"])
			} else {
				assert.Equal(suite.T(), "User retrieved successfully", response["message"])
				userData := response["data"].(map[string]interface{})
				assert.Equal(suite.T(), testUser.Email, userData["email"])
			}
		})
	}
}

func (suite *UserControllerTestSuite) TestUpdateUserController() {
	testUser := models.User{
		Name:  "Test User",
		Email: "test.update@example.com",
		Age:   30,
	}
	result := dao.GetDB().Create(&testUser)
	assert.NoError(suite.T(), result.Error)

	tests := []struct {
		name          string
		identifier    string
		updateData    models.User
		expectedCode  int
		expectedError string
	}{
		{
			name:       "Valid Update",
			identifier: "test.update@example.com",
			updateData: models.User{
				Name: "Updated Name",
				Age:  31,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:       "Invalid Age Update",
			identifier: "test.update@example.com",
			updateData: models.User{
				Age: 151,
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "Failed to update user: validation error: age must be less than 150",
		},
		{
			name:       "Update Non-existent User",
			identifier: "nonexistent@example.com",
			updateData: models.User{
				Name: "New Name",
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "Failed to update user: user not found",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			jsonBody, err := json.Marshal(tt.updateData)
			assert.NoError(suite.T(), err)

			req := httptest.NewRequest("PUT", "/v1/api/user/"+tt.identifier, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			web.BeeApp.Handlers.ServeHTTP(w, req)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Equal(suite.T(), tt.expectedError, response["error"])
			} else {
				assert.Equal(suite.T(), "User updated successfully", response["message"])
				if userData, ok := response["data"].(map[string]interface{}); ok {
					if tt.updateData.Name != "" {
						assert.Equal(suite.T(), tt.updateData.Name, userData["name"])
					}
					if tt.updateData.Age > 0 {
						assert.Equal(suite.T(), float64(tt.updateData.Age), userData["age"])
					}
				}
			}
		})
	}
}

func (suite *UserControllerTestSuite) TestDeleteUserController() {
	testUser := models.User{
		Name:  "Test User",
		Email: "test.delete@example.com",
		Age:   30,
	}
	result := dao.GetDB().Create(&testUser)
	assert.NoError(suite.T(), result.Error)

	tests := []struct {
		name          string
		identifier    string
		expectedCode  int
		expectedError string
	}{
		{
			name:         "Valid Delete",
			identifier:   "test.delete@example.com",
			expectedCode: http.StatusOK,
		},
		{
			name:          "Delete Non-existent User",
			identifier:    "nonexistent@example.com",
			expectedCode:  http.StatusNotFound,
			expectedError: "User not found: user not found",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req := httptest.NewRequest("DELETE", "/v1/api/user/"+tt.identifier, nil)
			w := httptest.NewRecorder()
			web.BeeApp.Handlers.ServeHTTP(w, req)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Equal(suite.T(), tt.expectedError, response["error"])
			} else {
				assert.Equal(suite.T(), "User deleted successfully", response["message"])
				userData := response["data"].(map[string]interface{})
				assert.Equal(suite.T(), testUser.Email, userData["email"])
			}
		})
	}
}