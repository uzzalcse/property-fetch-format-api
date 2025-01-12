package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"property-fetch-format-api/models"
	"regexp"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type PropertyDetailsController struct {
    beego.Controller
}

type PropertyService interface {
    FetchPropertyDetails(ctx context.Context, propertyID, languageCode string) (*models.PropertyResponse, error)
}

type propertyService struct {
    baseURL    string
    httpClient *http.Client
}

func NewPropertyService(baseURL string) PropertyService {
    return &propertyService{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (s *propertyService) FetchPropertyDetails(ctx context.Context, propertyID, languageCode string) (*models.PropertyResponse, error) {
    // Validate property ID format using regex
    if ok, err := regexp.MatchString(`^[A-Z]{2}-\d+$`, propertyID); !ok || err != nil {
        return nil, fmt.Errorf("invalid property ID format")
    }
    url := fmt.Sprintf("%s?propertyId=%s&languageCode=%s", s.baseURL, propertyID, languageCode)
    
    resultChan := make(chan *models.ExternalAPIResponse, 1)
    errChan := make(chan error, 1)
    
    go func() {
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        if err != nil {
            errChan <- fmt.Errorf("error creating request: %w", err)
            return
        }
        
        resp, err := s.httpClient.Do(req)
        if err != nil {
            errChan <- fmt.Errorf("error making request: %w", err)
            return
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            errChan <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
            return
        }
        
        var apiResp models.ExternalAPIResponse
        if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
            errChan <- fmt.Errorf("error decoding response: %w", err)
            return
        }
        
        resultChan <- &apiResp
    }()
    
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case err := <-errChan:
        return nil, err
    case result := <-resultChan:
        return transformToPropertyResponse(result), nil
    }
}

// @Title GetPropertyDetails
// @Description Get details of a property by its ID and language code
// @Param   propertyId    path      string  true  "Property ID in format 'XX-1234'"
// @Param   languageCode  query     string  false "Language code (default: en)"
// @Success 200 {object} models.PropertyResponse "Property details response"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "Property not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router  /v1/api/property/details/{propertyId} [get]
func (c *PropertyDetailsController) GetPropertyDetails() {
    ctx := c.Ctx.Request.Context()
    propertyID := c.Ctx.Input.Param(":propertyId")
    languageCode := c.GetString("languageCode", "en")
    
    if propertyID == "" {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        c.Data["json"] = map[string]string{"error": "property ID is required"}
        c.ServeJSON()
        return
    }
    
    baseURL, _ := beego.AppConfig.String("baseurl")
    svc := NewPropertyService(baseURL)
    
    property, err := svc.FetchPropertyDetails(ctx, propertyID, languageCode)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    if property == nil {
        c.Ctx.Output.SetStatus(http.StatusNotFound)
        c.Data["json"] = map[string]string{"error": "Property not found"}
        c.ServeJSON()
        return
    }
    
    c.Data["json"] = property
    c.ServeJSON()
}

func transformToPropertyResponse(apiResp *models.ExternalAPIResponse) *models.PropertyResponse {
    // Check if the response or critical fields are nil/empty
    if apiResp == nil || apiResp.S3.ID == "" {
        return nil
    }

    // Transform external API response to PropertyResponse
    return &models.PropertyResponse{
        ID:        apiResp.S3.ID,
        Feed:      apiResp.S3.Feed,
        Published: apiResp.S3.Published,
        GeoInfo:   apiResp.S3.GeoInfo,
        Property:  apiResp.S3.Property,
        Partner:   apiResp.S3.Partner,
    }
}