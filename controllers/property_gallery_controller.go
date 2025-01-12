package controllers

import (
	"context"
	"encoding/json"
	"fmt"
    "property-fetch-format-api/models"
	"net/http"
	"regexp"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type PropertyGalleryController struct {
    beego.Controller
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
    
    resultChan := make(chan *models.GalleryResponse, 1)
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

        var gallery models.GalleryResponse
        if err := json.NewDecoder(resp.Body).Decode(&gallery); err != nil {
            errChan <- fmt.Errorf("property not found")
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

func transformToGroupedImages(gallery *models.GalleryResponse) GroupedImages {
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

// @Summary Get property gallery images
// @Description Retrieve property gallery images grouped by labels.
// @Tags Property Gallery
// @Param propertyId path string true "Property ID in format XX-123"
// @Param languageCode query string false "Language code for the images" default(en)
// @Success 200 {object} GroupedImages "Grouped images by label"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/api/property/{propertyId}/gallery [get]
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
    if ok, err := regexp.MatchString(`^[A-Z]{2}-\d+$`, propertyID); !ok || err != nil {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        c.Data["json"] = map[string]string{"error": "invalid property ID format"}
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