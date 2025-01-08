package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "property-fetch-format-api/models"
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
    
    c.Data["json"] = property
    c.ServeJSON()
}

func transformToPropertyResponse(apiResp *models.ExternalAPIResponse) *models.PropertyResponse {
    // Transform external API response to PropertyResponse
    return &models.PropertyResponse{
        ID:        apiResp.S3.ID,
        Feed:      apiResp.S3.Feed, // Example value
        Published: apiResp.S3.Published,
        GeoInfo:   apiResp.S3.GeoInfo,
        Property:  apiResp.S3.Property,
        Partner:   apiResp.S3.Partner,
    }
}