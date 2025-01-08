package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyGalleryController struct {
    beego.Controller
}

type GalleryImage struct {
    Captions    string      `json:"captions"`
    Confidence  float64     `json:"confidence"`
    Height      int         `json:"height"`
    ID          int         `json:"id"`
    Label       string      `json:"label"`
    Predictions interface{} `json:"predictions"`
    URL         string      `json:"url"`
}

type GalleryResponse struct {
    S3Gallery map[string][]GalleryImage `json:"S3-Gallery"`
}

type GroupedImages map[string][]string

type GalleryService interface {
    FetchPropertyGallery(ctx context.Context, propertyID, languageCode string) (GroupedImages, error)
}

type galleryService struct {
    baseURL    string
    httpClient *http.Client
}

func NewGalleryService(baseURL string) GalleryService {
    return &galleryService{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (s *galleryService) FetchPropertyGallery(ctx context.Context, propertyID, languageCode string) (GroupedImages, error) {
    url := fmt.Sprintf("%s?propertyId=%s&languageCode=%s", s.baseURL, propertyID, languageCode)
    
    resultChan := make(chan *GalleryResponse, 1)
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
        
        var gallery GalleryResponse
        if err := json.NewDecoder(resp.Body).Decode(&gallery); err != nil {
            errChan <- fmt.Errorf("error decoding response: %w", err)
            return
        }
        
        resultChan <- &gallery
    }()
    
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case err := <-errChan:
        return nil, err
    case result := <-resultChan:
        return transformToGroupedImages(result), nil
    }
}

func transformToGroupedImages(gallery *GalleryResponse) GroupedImages {
    grouped := GroupedImages{}
    
    for _, images := range gallery.S3Gallery {
        for _, img := range images {
            if img.Confidence > 95 {
                if _, exists := grouped[img.Label]; !exists {
                    grouped[img.Label] = []string{}
                }
                grouped[img.Label] = append(grouped[img.Label], img.URL)
            }
        }
    }
    
    return grouped
}

func (c *PropertyGalleryController) GetPropertyGallery() {
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
    svc := NewGalleryService(baseURL)
    
    response, err := svc.FetchPropertyGallery(ctx, propertyID, languageCode)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    c.Data["json"] = response
    c.ServeJSON()
}