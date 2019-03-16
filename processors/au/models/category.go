package aumodels

import "time"

// Category is the root element of a category output
type Category struct {
    Adverts                 []NormalisedAdvert  `json:"ads"`
    MatchedEntries          uint                `json:"matched_entries"`
    Pagination              *Pagination         `json:"pagination"`
}

// Pagination is the page control of the category
type Pagination struct {
    CurrentPage             uint                `json:"current"`
    PageSize                uint                `json:"size"`
}

// NormalisedAdvert is the root element of an ad
type NormalisedAdvert struct {
    ID                      uint                `json:"id"`
    Type                    string              `json:"type"`
    Status                  string              `json:"status"`
    Category                *AdvertCategory     `json:"category"`
    Position                *Position           `json:"positions"`
    PosterType              *string             `json:"poster_type"`
    Price                   *Price              `json:"price"`
    Title                   string              `json:"title"`
    DescriptionExcerpt      *string             `json:"desc_excerpt"`
    DescriptionExcerptHTML  *string             `json:"desc_excerpt_html"`
    Pictures                []Picture           `json:"pictures"`
    Attributes              []Attribute         `json:"attributes"`
    CreationTime            time.Time           `json:"creation_time"`
    StartTime               time.Time           `json:"start_time"`
    EndTime                 time.Time           `json:"end_time"`
}

// Price is the listing price shown on an ad
type Price struct {
    Type                    *string             `json:"type"`
    Amount                  *uint               `json:"amount"`
    HighestAmount           *uint               `json:"highest_amount"`
    Currency                *string             `json:"currency"`
    CurrencySymbol          *string             `json:"currency_symbol"`
}

// Position is the positional information of an ad
type Position struct {
    Coordinate              *Coordinate         `json:"coordinate"`
    State                   *string             `json:"state"`
    Locations               []Location          `json:"locations"`
}

// Location is the locality information of an ad
type Location struct {
    ID                      uint                `json:"id"`
    Name                    string              `json:"name"`
    Regions                 []Region            `json:"regions"`
}

// Region is list of address levels inside the location
type Region struct {
    Level                   string              `json:"level"`
    Name                    string              `json:"name"`
}

// Coordinate is the geographic coordinates of the location
type Coordinate struct {
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

// Attribute associated with each ad
type Attribute struct {
    KeySlug                 string              `json:"key_slug"`
    KeyName                 *string             `json:"key_name"`
    ValueSlug               *string             `json:"value_slug"`
    ValueName               *string             `json:"value_name"`
}

// Picture associated with each ad
type Picture struct {
    Thumbnail               *string             `json:"thumbnail_url"`
    Normal                  *string             `json:"normal_url"`
    Large                   *string             `json:"large_url"`
    ExtraLarge              *string             `json:"extra_large"`
    ExtraExtraLarge         *string             `json:"extra_extra_large"`
}
