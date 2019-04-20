package auparser

import (
    "fmt"
    models "github.com/GreenVine/ebay-classifieds-api/parsers/au/models"
    u "github.com/GreenVine/ebay-classifieds-api/utils"
    "github.com/beevik/etree"
    "strings"
    "time"
)

// ParseAdvert is to build a Category models from raw XML response
func ParseAdvert(doc *etree.Document) (*models.Advert, []error, bool) {
    if doc == nil {
        return nil, []error{ fmt.Errorf("empty API response") }, true
    }

    root := doc.Root()

    if root == nil || root.Space != "ad" || root.Tag != "ad" {
        return nil, []error{ fmt.Errorf("unexpected API response") }, true
    }

    var errors []error
    var hasCriticalError = false

    if advert := BuildAdvertBase(root, &errors, &hasCriticalError); !hasCriticalError {
        return advert, errors, false
    }

    return nil, errors, true
}

// BuildAdvertBase is to build the base of an advertisement
func BuildAdvertBase(ad *etree.Element, errors *[]error, hasCriticalError *bool) *models.Advert {
    advertID, err := u.ConvString2Uint(u.ExtractAttrByTag(ad, "id"))
    if err != nil {
        *errors = append(*errors, fmt.Errorf("ads/ad/id"))
        return nil
    }

    advertType, _ := u.ExtractText(ad, "./ad:ad-type/ad:value")

    advertUserID, _ := u.ConvString2Uint(u.ExtractText(ad, "./ad:user-id"))

    advertPrice := buildPrice(ad, errors, hasCriticalError)

    advertStatus, _ := u.ExtractText(ad, "./ad:ad-status/ad:value")

    advertContact := buildContact(ad, errors, hasCriticalError)

    advertCategory := buildCategory(ad, errors, hasCriticalError)

    advertPosition := buildPosition(ad, errors, hasCriticalError)

    advertPosterType, _ := u.ExtractText(ad, "./ad:poster-type/ad:value")

    advertTitle := u.FallbackStringWithReport(
        u.ExtractText(ad, "./ad:title"))(
        "", errors, fmt.Errorf("ads/ad/title"))

    advertDescriptionExcerptHTML := u.FallbackStringWithReport(
        u.ExtractText(ad, "./ad:description"))(
        "", errors, fmt.Errorf("ads/ad/desc_excerpt_html"))

    advertDescriptionExcerpt, _ := u.FormatHTML2Base64(advertDescriptionExcerptHTML)

    advertPictures := buildPicture(ad, errors, hasCriticalError)

    advertAttributes := buildAttribute(ad, errors, hasCriticalError)
    advertTimestamp := buildTimestamp(ad, errors, hasCriticalError)

    return &models.Advert{
        ID:         advertID,
        Type:       u.ReplaceStringWithNil(&advertType, ""),
        UserID:     u.ReplaceUintWithNil(&advertUserID, 0),
        Status:     u.ReplaceStringWithNil(&advertStatus, ""),
        Contact:    advertContact,
        Category:   advertCategory,
        Position:   advertPosition,
        PosterType: u.ReplaceStringWithNil(&advertPosterType, ""),
        Price:      advertPrice,
        Title:      *u.ReplaceStringWithNil(&advertTitle, ""),
        DescriptionExcerptB64:  u.ReplaceStringWithNil(&advertDescriptionExcerpt, ""),
        DescriptionExcerptHTML: u.ReplaceStringWithNil(&advertDescriptionExcerptHTML, ""),
        Pictures:               advertPictures,
        Attributes:             advertAttributes,
        Timestamp:              advertTimestamp,
    }
}

func buildPrice(ad *etree.Element, errors *[]error, _ *bool) *models.AdvertPrice {
    priceType := u.FallbackStringWithReport(
        u.ExtractText(ad, "./ad:price/types:price-type/types:value"))(
        "UNKNOWN", errors, fmt.Errorf("ads/ad/price/type"))

    priceAmount := uint(u.FallbackFloat64WithReport(
        u.ExtractTextAsFloat64(ad, "./ad:price/types:amount"))(
        0.0, errors, fmt.Errorf("ads/ad/price/amount")) * 100)

    priceHighestAmount := uint(u.FallbackFloat64WithReport(
        u.ExtractTextAsFloat64(ad, "./ad:highest-price"))(
        0, errors, fmt.Errorf("ads/ad/price/highest_amount")) * 100)

    currency := u.FallbackStringWithReport(
        u.ExtractText(ad, "./ad:price/types:currency-iso-code/types:value"))(
        "", errors, fmt.Errorf("ads/ad/price/currency"))

    currencySymbol := u.FallbackStringWithReport(
        u.ExtractAttrByTag(ad.FindElement(
            "./ad:price/types:currency-iso-code/types:value"), "localized-label"))(
        "", errors, fmt.Errorf("ads/ad/price/currency_symbol"))

    return &models.AdvertPrice{
        Type: &priceType,
        Amount: &priceAmount,
        HighestAmount: &priceHighestAmount,
        Currency: u.ReplaceStringWithNil(&currency, ""),
        CurrencySymbol: u.ReplaceStringWithNil(&currencySymbol, ""),
    }
}

func buildCategory(ad *etree.Element, errors *[]error, _ *bool) *models.AdvertCategory {
    cat := ad.FindElement("./cat:category")
    if cat == nil {
        return nil
    }

    catID, err := u.ConvString2Uint(u.ExtractAttrByTag(cat, "id"))
    if err != nil {
        return nil
    }

    catName := u.FallbackStringWithReport(
        u.ExtractText(cat, "./cat:localized-name"))(
        "", errors, fmt.Errorf("ads/ad/category/name"))

    catSlug := u.FallbackStringWithReport(
        u.ExtractText(cat, "./cat:id-name"))(
        "", errors, fmt.Errorf("ads/ad/category/slug"))

    catParentSlug := u.FallbackStringWithReport(
        u.ExtractText(cat, "./cat:l1-name"))(
        "", errors, fmt.Errorf("ads/ad/category/parent_slug"))

    catChildrenCount := u.FallbackUintWithReport(
        u.ExtractTextAsUint(cat, "./cat:children-count"))(
        0, errors, fmt.Errorf("ads/ad/category/children_count"))

    return &models.AdvertCategory{
        ID:            catID,
        Name:          *u.ReplaceStringWithNil(&catName, ""),
        Slug:          u.ReplaceStringWithNil(&catSlug, ""),
        ParentSlug:    u.ReplaceStringWithNil(&catParentSlug, ""),
        ChildrenCount: &catChildrenCount,
    }
}

func buildPosition(ad *etree.Element, errors *[]error, _ *bool) *models.AdvertPosition {
    var coordinate *models.AdvertCoordinate
    var locations []models.AdvertLocation

    address, _ := u.ExtractText(ad, "./ad:ad-address/types:full-address")
    city, _ := u.ExtractText(ad, "./ad:ad-address/types:city")
    state, _ := u.ExtractText(ad, "./ad:ad-address/types:state")
    country, _ := u.ExtractText(ad, "./ad:ad-address/types:country")

    longitude, longerr := u.ConvString2Float64(u.ExtractText(ad, "./ad:ad-address/types:longitude"))
    latitude, laterr := u.ConvString2Float64(u.ExtractText(ad, "./ad:ad-address/types:latitude"))

    if longerr == nil && laterr == nil {
        coordinate = &models.AdvertCoordinate{
            Longitude: longitude,
            Latitude: latitude,
        }
    } else {
        *errors = append(*errors, fmt.Errorf("ads/ad/positions/coordinate"))
    }

    if locs := ad.FindElements("./loc:locations/loc:location"); locs != nil {
        for i, loc := range locs {
            locID := u.FallbackUintWithReport(
                u.ConvString2Uint(u.ExtractAttrByTag(loc, "id")))(
                0, errors, fmt.Errorf("ads/ad/positions/locations[%d]/id", i))

            locName := u.FallbackStringWithReport(
                u.ExtractText(loc, "./loc:localized-name"))(
                "", errors, fmt.Errorf("ads/ad/positions/location[%d]/name", i))

            locParentID, _ := u.ConvString2Uint(u.ExtractText(loc, "./loc:parent-id"))

            locations = append(locations, models.AdvertLocation{
                ID:         *u.ReplaceUintWithNil(&locID, 0),
                Name:       *u.ReplaceStringWithNil(&locName, ""),
                ParentID:   u.ReplaceUintWithNil(&locParentID, 0),
            })
        }
    }

    return &models.AdvertPosition{
        Address:    u.ReplaceStringWithNil(&address, ""),
        City:       u.ReplaceStringWithNil(&city, ""),
        State:      u.ReplaceStringWithNil(&state, ""),
        Country:    u.ReplaceStringWithNil(&country, ""),
        Coordinate: coordinate,
        Locations:  locations,
    }
}

func buildPicture(ad *etree.Element, errors *[]error, _ *bool) []models.AdvertPicture {
    var pictures []models.AdvertPicture

    if pics := ad.FindElements("./pic:pictures/pic:picture"); pics != nil {
        for i, pic := range pics {
            if pic != nil {
                thumbnail := u.FallbackStringWithReport(
                    u.ExtractAttrByTag(pic.FindElement("./pic:link[@rel='thumbnail']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/thumbnail", i))
                normal := u.FallbackStringWithReport(
                    u.ExtractAttrByTag(pic.FindElement("./pic:link[@rel='normal']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/normal", i))
                large := u.FallbackStringWithReport(
                    u.ExtractAttrByTag(pic.FindElement("./pic:link[@rel='large']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/large", i))
                extraLarge := u.FallbackStringWithReport(
                    u.ExtractAttrByTag(pic.FindElement("./pic:link[@rel='extraLarge']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/extraLarge", i))
                extra2XLarge := u.FallbackStringWithReport(
                    u.ExtractAttrByTag(pic.FindElement("./pic:link[@rel='extraExtraLarge']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/extra2XLarge", i))

                pictures = append(pictures, models.AdvertPicture{
                    Thumbnail:    u.ReplaceStringWithNil(&thumbnail, ""),
                    Normal:       u.ReplaceStringWithNil(&normal, ""),
                    Large:        u.ReplaceStringWithNil(&large, ""),
                    ExtraLarge:   u.ReplaceStringWithNil(&extraLarge, ""),
                    Extra2XLarge: u.ReplaceStringWithNil(&extra2XLarge, ""),
                })
            }
        }
    }

    return pictures
}

func buildAttribute(ad *etree.Element, errors *[]error, _ *bool) []models.AdvertAttribute {
    var attributes []models.AdvertAttribute

    if attrs := ad.FindElements("./attr:attributes/attr:attribute"); attrs != nil {
        for i, attr := range attrs {
            if attr != nil {
                keySlug := u.FallbackStringWithReport(
                    u.ExtractAttrByTag(attr, "name"))(
                    "", errors, fmt.Errorf("ads/ad/attributes[%d]/key_slug", i))
                keyName := u.FallbackStringWithReport(
                    u.ExtractAttrByTag(attr, "localized-label"))(
                    "", errors, fmt.Errorf("ads/ad/attributes[%d]/key_name", i))
                valueType := u.FallbackStringWithReport(
                    u.ExtractAttrByTag(attr, "type"))(
                    "", errors, fmt.Errorf("ads/ad/attributes[%d]/value_type", i))

                valueSlugStr, _ := u.ExtractText(attr, "./attr:value")
                valueSlug := u.ReplaceStringWithNil(&valueSlugStr, "")

                valueNameStr, _ := u.ExtractAttrByTag(attr.FindElement("./attr:value"), "localized-label")
                valueName := u.ReplaceStringWithNil(&valueNameStr, "")

                attributes = append(attributes, models.AdvertAttribute{
                    KeySlug:    *u.ReplaceStringWithNil(&keySlug, ""),
                    KeyName:    *u.ReplaceStringWithNil(&keyName, ""),
                    ValueType:  u.ReplaceStringWithNil(&valueType, ""),
                    ValueSlug:  u.ReplaceStringWithNil(valueSlug, ""),
                    ValueName:  u.ReplaceStringWithNil(valueName, ""),
                })
            }
        }
    }

    return attributes
}

func buildTimestamp(ad *etree.Element, _ *[]error, _ *bool) models.AdvertTimestamp {
    advertCreationTime := formatTimestamp(ad.FindElement("./ad:creation-date-time"))
    advertModificationTime := formatTimestamp(ad.FindElement("./ad:modification-date-time"))
    advertStartTime := formatTimestamp(ad.FindElement("./ad:start-date-time"))
    advertEndTime := formatTimestamp(ad.FindElement("./ad:end-date-time"))

    return models.AdvertTimestamp{
        CreationTime:       advertCreationTime,
        ModificationTime:   advertModificationTime,
        StartTime:          advertStartTime,
        EndTime:            advertEndTime,
    }
}

func buildContact(ad *etree.Element, _ *[]error, _ *bool) *models.AdvertContact {
    name, _  := u.ExtractText(ad, "./ad:poster-contact-name")
    phone, _ := u.ExtractText(ad, "./ad:phone")
    phone     = strings.Replace(phone, " ", "", -1)

    return &models.AdvertContact{
        Name:   u.ReplaceStringWithNil(&name, ""),
        Phone:  u.ReplaceStringWithNil(&phone, ""),
    }
}

func formatTimestamp(element *etree.Element) *time.Time {
    if element != nil {
        timestr := element.Text()

        if timeinst, err := time.Parse(time.RFC3339, timestr); err == nil {
            timeinst = timeinst.UTC()
            return &timeinst
        }
    }

    return nil
}
