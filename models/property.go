// models/property.go
package models

type PropertyResponse struct {
    ID        string    `json:"ID"`
    Feed      int       `json:"Feed"`
    Published bool      `json:"Published"`
    GeoInfo   GeoInfo   `json:"GeoInfo"`
    Property  Property  `json:"Property"`
    Partner   Partner   `json:"Partner"`
}

type GeoInfo struct {
    Categories   []Category `json:"Categories"`
    City         string     `json:"City"`
    Country      string     `json:"Country"`
    CountryCode  string     `json:"CountryCode"`
    Display      string     `json:"Display"`
    LocationID   string     `json:"LocationID"`
    StateAbbr    string     `json:"StateAbbr"`
    Lat          string     `json:"Lat"`
    Lng          string     `json:"Lng"`
}

type Category struct {
    Name       string   `json:"Name"`
    Slug       string   `json:"Slug"`
    Type       string   `json:"Type"`
    Display    []string `json:"Display"`
    LocationID string   `json:"LocationID"`
}

type Property struct {
    Amenities    map[string]string `json:"Amenities"`
    Counts       Counts            `json:"Counts"`
    EcoFriendly  bool              `json:"EcoFriendly"`
    FeatureImage string            `json:"FeatureImage"`
    Image        Image             `json:"Image"`
    Price        float64           `json:"Price"`
    PropertyName string            `json:"PropertyName"`
    PropertySlug string            `json:"PropertySlug"`
    PropertyType string            `json:"PropertyType"`
	PropertyTypeCategoryId string	 `json:"PropertyTypeCategoryId"`
    ReviewScore  float64           `json:"ReviewScore"`
    RoomSize     int              `json:"RoomSize"`
    MinStay      int              `json:"MinStay"`
    UpdatedAt    string           `json:"UpdatedAt"`
}

type Counts struct {
    Bedroom   int `json:"Bedroom"`
    Bathroom  int `json:"Bathroom"`
    Reviews   int `json:"Reviews"`
    Occupancy int `json:"Occupancy"`
}

type Image struct {
    Count  int      `json:"Count"`
    Images []string `json:"Images"`
}

type Partner struct {
    ID       string   `json:"ID"`
    Archived []string `json:"Archived"`
    OwnerID  string   `json:"OwnerID"`
    HcomID   string   `json:"HcomID"`
    BrandId  string   `json:"BrandId"`
	 URL      string   `json:"URL"`
	UnitNumber string `json:"UnitNumber"`
	EpCluster  string `json:"EpCluster"`
}

// ExternalAPIResponse represents the combined response from external API
type ExternalAPIResponse struct {
    S3 S3Response `json:"S3"`
}

type S3Response struct {
    // Add relevant fields from S3 response
    ID        string    `json:"ID"`
	Feed      int       `json:"Feed"`
    Property  Property  `json:"Property"`
    Partner   Partner   `json:"Partner"`
    GeoInfo   GeoInfo   `json:"GeoInfo"`
    Published bool      `json:"Published"`
}
