package aumodels

import "time"

// Advert is the root element of an ad
type Advert struct {
    ID                      uint                `json:"id"`
    Type                    *string             `json:"type"`
    UserID                  *uint               `json:"user_id"`
    Status                  *string             `json:"status"`
    Contact                 *AdvertContact      `json:"contact"`
    Category                *AdvertCategory     `json:"category"`
    Position                *AdvertPosition     `json:"positions"`
    PosterType              *string             `json:"poster_type"`
    Price                   *AdvertPrice        `json:"price"`
    Title                   string              `json:"title"`
    DescriptionExcerptB64   *string             `json:"desc_excerpt_plain_b64"`
    DescriptionExcerptHTML  *string             `json:"desc_excerpt_html"`
    Pictures                []AdvertPicture     `json:"pictures"`
    Attributes              []AdvertAttribute   `json:"attributes"`
    Timestamp               AdvertTimestamp     `json:"timestamp"`
}

// AdvertPrice is the listing price shown on an ad
type AdvertPrice struct {
    Type                    *string             `json:"type"`
    Amount                  *uint               `json:"amount"`
    HighestAmount           *uint               `json:"highest_amount"`
    Currency                *string             `json:"currency"`
    CurrencySymbol          *string             `json:"currency_symbol"`
}

// AdvertPosition is the positional information of an ad
type AdvertPosition struct {
    Address                 *string             `json:"address"`
    City                    *string             `json:"city"`
    State                   *string             `json:"state"`
    Country                 *string             `json:"country"`
    Coordinate              *AdvertCoordinate   `json:"coordinate"`
    Locations               []AdvertLocation    `json:"locations"`
}

// AdvertLocation is the locality information of an ad
type AdvertLocation struct {
    ID                      uint                `json:"id"`
    Name                    string              `json:"name"`
    ParentID                *uint               `json:"parent_id"`
}

// AdvertCoordinate is the geographic coordinates of the location
type AdvertCoordinate struct {
    Longitude               float64             `json:"longitude"`
    Latitude                float64             `json:"latitude"`
}

// AdvertCategory is the category of an ad (not to be confused with the root category)
type AdvertCategory struct {
    ID                      uint                `json:"id"`
    Name                    string              `json:"name"`
    Slug                    *string             `json:"slug"`
    ParentSlug              *string             `json:"parent_slug"`
    ChildrenCount           *uint               `json:"children_count"`
}

// AdvertAttribute is the attribute associated with an ad
type AdvertAttribute struct {
    KeySlug                 string              `json:"key_slug"`
    KeyName                 string              `json:"key_name"`
    ValueType               *string             `json:"value_type"`
    ValueSlug               *string             `json:"value_slug"`
    ValueName               *string             `json:"value_name"`
}

// AdvertPicture is the picture associated with an ad
type AdvertPicture struct {
    Thumbnail               *string             `json:"thumbnail_url"`
    Normal                  *string             `json:"normal_url"`
    Large                   *string             `json:"large_url"`
    ExtraLarge              *string             `json:"extra_large_url"`
    Extra2XLarge            *string             `json:"extra_2x_large_url"`
}

// AdvertTimestamp is the timestamp associated with an ad
type AdvertTimestamp struct {
    CreationTime            *time.Time          `json:"creation_time"`
    ModificationTime        *time.Time          `json:"modification_time"`
    StartTime               *time.Time          `json:"start_time"`
    EndTime                 *time.Time          `json:"end_time"`
}

// AdvertContact is the contact details of an ad
type AdvertContact struct {
    Name                    *string             `json:"name"`
    Phone                   *string             `json:"phone"`
}
