package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "sync"
    "time"
    "property-fetch-format-api/models"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyListController struct {
    beego.Controller
}

// PropertyListService interface defines bulk property fetching operations
type PropertyListService interface {
    FetchPropertyList(ctx context.Context, propertyIDs []string, languageCode string) ([]models.PropertyResponse, error)
}

type propertyListService struct {
    baseURL    string
    httpClient *http.Client
}

// NewPropertyListService creates a new instance of PropertyListService
func NewPropertyListService(baseURL string) PropertyListService {
    return &propertyListService{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// FetchPropertyList fetches multiple properties in parallel
func (s *propertyListService) FetchPropertyList(ctx context.Context, propertyIDs []string, languageCode string) ([]models.PropertyResponse, error) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    results := make([]models.PropertyResponse, len(propertyIDs))
    errors := make(chan error, len(propertyIDs))

    for i, id := range propertyIDs {
        wg.Add(1)
        go func(index int, propertyID string) {
            defer wg.Done()

            url := fmt.Sprintf("%s?propertyId=%s&languageCode=%s", s.baseURL, propertyID, languageCode)
            
            req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
            if err != nil {
                errors <- fmt.Errorf("error creating request for property %s: %w", propertyID, err)
                return
            }

            resp, err := s.httpClient.Do(req)
            if err != nil {
                errors <- fmt.Errorf("error fetching property %s: %w", propertyID, err)
                return
            }
            defer resp.Body.Close()

            if resp.StatusCode != http.StatusOK {
                errors <- fmt.Errorf("unexpected status code %d for property %s", resp.StatusCode, propertyID)
                return
            }

            var apiResp models.ExternalAPIResponse
            if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
                errors <- fmt.Errorf("error decoding response for property %s: %w", propertyID, err)
                return
            }

            transformedData := transformToPropertyResponse(&apiResp)

            mu.Lock()
            results[index] = *transformedData
            mu.Unlock()
        }(i, id)
    }

    // Create error collection channel
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    // Wait for either completion or context cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case <-done:
        // Check if there were any errors
        select {
        case err := <-errors:
            return nil, err
        default:
            return results, nil
        }
    }
}

// GetPropertyList handles the HTTP request for bulk property fetching
func (c *PropertyListController) GetPropertyList() {
    ctx := c.Ctx.Request.Context()
    propertyIDs := strings.Split(c.GetString("propertyIds"), ",")
    languageCode := c.GetString("languageCode", "en")

    if len(propertyIDs) == 0 {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        c.Data["json"] = map[string]string{"error": "property IDs are required"}
        c.ServeJSON()
        return
    }

    baseURL, _ := beego.AppConfig.String("baseurl")
    svc := NewPropertyListService(baseURL)

    response, err := svc.FetchPropertyList(ctx, propertyIDs, languageCode)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    c.Data["json"] = response
    c.ServeJSON()
}