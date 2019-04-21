package aumodels

import "time"

// Advert is the root element of an ad
type Advert struct {
    ID                      uint                `json:"id"`
    Type                    *string             `json:"type"`
    UserID                  *uint               `json:"user_id"`
    Status                  *string             `json:"status"`
    Contact                 *AdvertContact      `json:"contact"`
    Category                *AdvertCategory     `json:"category,omitempty"`
    Position                *AdvertPosition     `json:"positions,omitempty"`
    PosterType              *string             `json:"poster_type,omitempty"`
    Price                   *AdvertPrice        `json:"price"`
    Title                   string              `json:"title"`
    DescriptionExcerptB64   *string             `json:"desc_excerpt_plain_b64,omitempty"`
    DescriptionExcerptHTML  *string             `json:"desc_excerpt_html,omitempty"`
    Pictures                []AdvertPicture     `json:"pictures,omitempty"`
    Attributes              []AdvertAttribute   `json:"attributes,omitempty"`
    Timestamp               AdvertTimestamp     `json:"timestamp"`
}

// AdvertPrice is the listing price shown on an ad
type AdvertPrice struct {
    Type                    *string             `json:"type"`
    Amount                  *uint               `json:"amount"`
    HighestAmount           *uint               `json:"highest_amount,omitempty"`
    Currency                *string             `json:"currency,omitempty"`
    CurrencySymbol          *string             `json:"currency_symbol,omitempty"`
}

// AdvertPosition is the positional information of an ad
type AdvertPosition struct {
    Address                 *string             `json:"address,omitempty"`
    City                    *string             `json:"city"`
    State                   *string             `json:"state"`
    Country                 *string             `json:"country"`
    Coordinate              *AdvertCoordinate   `json:"coordinate,omitempty"`
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
    Slug                    *string             `json:"slug,omitempty"`
    ParentSlug              *string             `json:"parent_slug,omitempty"`
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
    Thumbnail               *string             `json:"thumbnail_url,omitempty"`
    Normal                  *string             `json:"normal_url,omitempty"`
    Large                   *string             `json:"large_url,omitempty"`
    ExtraLarge              *string             `json:"extra_large_url,omitempty"`
    Extra2XLarge            *string             `json:"extra_2x_large_url,omitempty"`
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
    Phone                   *string             `json:"phone,omitempty"`
}
