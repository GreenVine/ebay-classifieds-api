package auparsers

import (
    "fmt"
    . "github.com/GreenVine/ebay-ecg-api/processors/au/models"
    . "github.com/GreenVine/ebay-ecg-api/utils"
    "github.com/beevik/etree"
    "strings"
    "time"
)

// ParseAdvert is to build a Category models from raw XML response
func ParseAdvert(doc *etree.Document) (*Advert, []error, bool) {
    root := doc.Root()

    if root == nil || root.Space != "ad" || root.Tag != "ad" {
        return nil, []error{ fmt.Errorf("unexpected API response") }, true
    }

    var errors []error
    var hasCriticalError = false

    if advert := BuildAdvertBase(root, &errors, &hasCriticalError); hasCriticalError {
        return nil, errors, true
    } else {
        return advert, errors, false
    }
}

func BuildAdvertBase(ad *etree.Element, errors *[]error, hasCriticalError *bool) *Advert {
    advertId, err := ConvString2Uint(ExtractAttrByTag(ad, "id"))
    if err != nil {
        *errors = append(*errors, fmt.Errorf("ads/ad/id"))
        return nil
    }

    advertType, _ := ExtractText(ad, "./ad:ad-type/ad:value")

    advertUserId, _ := ConvString2Uint(ExtractText(ad, "./ad:user-id"))

    advertPrice := buildPrice(ad, errors, hasCriticalError)

    advertStatus, _ := ExtractText(ad, "./ad:ad-status/ad:value")

    advertContact := buildContact(ad, errors, hasCriticalError)

    advertCategory := buildCategory(ad, errors, hasCriticalError)

    advertPosition := buildPosition(ad, errors, hasCriticalError)

    advertPosterType, _ := ExtractText(ad, "./ad:poster-type/ad:value")

    advertTitle := FallbackStringWithReport(
        ExtractText(ad, "./ad:title"))(
        "", errors, fmt.Errorf("ads/ad/title"))

    advertDescriptionExcerptHTML := FallbackStringWithReport(
        ExtractText(ad, "./ad:description"))(
        "", errors, fmt.Errorf("ads/ad/desc_excerpt_html"))

    advertDescriptionExcerpt, _ := FormatHtml2Base64(advertDescriptionExcerptHTML)

    advertPictures := buildPicture(ad, errors, hasCriticalError)

    advertAttributes := buildAttribute(ad, errors, hasCriticalError)

    advertTimestamp := buildTimestamp(ad, errors, hasCriticalError)

    return &Advert{
        ID:                     advertId,
        Type:                   ReplaceStringWithNil(&advertType, ""),
        UserID:                 ReplaceUintWithNil(&advertUserId, 0),
        Status:                 ReplaceStringWithNil(&advertStatus, ""),
        Contact:                advertContact,
        Category:               advertCategory,
        Position:               advertPosition,
        PosterType:             ReplaceStringWithNil(&advertPosterType, ""),
        Price:                  advertPrice,
        Title:                  *ReplaceStringWithNil(&advertTitle, ""),
        DescriptionExcerptB64:  ReplaceStringWithNil(&advertDescriptionExcerpt, ""),
        DescriptionExcerptHTML: ReplaceStringWithNil(&advertDescriptionExcerptHTML, ""),
        Pictures:               advertPictures,
        Attributes:             advertAttributes,
        Timestamp:              advertTimestamp,
    }
}

func buildPrice(ad *etree.Element, errors *[]error, _ *bool) *AdvertPrice {
    priceType := FallbackStringWithReport(
        ExtractText(ad, "./ad:price/types:price-type/types:value"))(
        "UNKNOWN", errors, fmt.Errorf("ads/ad/price/type"))

    priceAmount := uint(FallbackFloat64WithReport(
        ExtractTextAsFloat64(ad, "./ad:price/types:amount"))(
        0.0, errors, fmt.Errorf("ads/ad/price/amount")) * 100)

    priceHighestAmount := uint(FallbackFloat64WithReport(
        ExtractTextAsFloat64(ad, "./ad:highest-price"))(
        0, errors, fmt.Errorf("ads/ad/price/highest_amount")) * 100)

    currency := FallbackStringWithReport(
        ExtractText(ad, "./ad:price/types:currency-iso-code/types:value"))(
        "", errors, fmt.Errorf("ads/ad/price/currency"))

    currencySymbol := FallbackStringWithReport(
        ExtractAttrByTag(ad.FindElement(
            "./ad:price/types:currency-iso-code/types:value"), "localized-label"))(
        "", errors, fmt.Errorf("ads/ad/price/currency_symbol"))

    return &AdvertPrice{
        Type: &priceType,
        Amount: &priceAmount,
        HighestAmount: &priceHighestAmount,
        Currency: ReplaceStringWithNil(&currency, ""),
        CurrencySymbol: ReplaceStringWithNil(&currencySymbol, ""),
    }
}

func buildCategory(ad *etree.Element, errors *[]error, _ *bool) *AdvertCategory {
    cat := ad.FindElement("./cat:category")
    if cat == nil {
        return nil
    }

    catId, err := ConvString2Uint(ExtractAttrByTag(cat, "id"))
    if err != nil {
        return nil
    }

    catName := FallbackStringWithReport(
        ExtractText(cat, "./cat:localized-name"))(
        "", errors, fmt.Errorf("ads/ad/category/name"))

    catSlug := FallbackStringWithReport(
        ExtractText(cat, "./cat:id-name"))(
        "", errors, fmt.Errorf("ads/ad/category/slug"))

    catParentSlug := FallbackStringWithReport(
        ExtractText(cat, "./cat:l1-name"))(
        "", errors, fmt.Errorf("ads/ad/category/parent_slug"))

    catChildrenCount := FallbackUintWithReport(
        ExtractTextAsUint(cat, "./cat:children-count"))(
        0, errors, fmt.Errorf("ads/ad/category/children_count"))

    return &AdvertCategory{
        ID:             catId,
        Name:           *ReplaceStringWithNil(&catName, ""),
        Slug:           ReplaceStringWithNil(&catSlug, ""),
        ParentSlug:     ReplaceStringWithNil(&catParentSlug, ""),
        ChildrenCount:  &catChildrenCount,
    }
}

func buildPosition(ad *etree.Element, errors *[]error, _ *bool) *AdvertPosition {
    var coordinate *AdvertCoordinate
    var locations []AdvertLocation

    address, _ := ExtractText(ad, "./ad:ad-address/types:full-address")
    city, _ := ExtractText(ad, "./ad:ad-address/types:city")
    state, _ := ExtractText(ad, "./ad:ad-address/types:state")
    country, _ := ExtractText(ad, "./ad:ad-address/types:country")

    longitude, longerr := ConvString2Float64(ExtractText(ad, "./ad:ad-address/types:longitude"))
    latitude, laterr := ConvString2Float64(ExtractText(ad, "./ad:ad-address/types:latitude"))

    if longerr == nil && laterr == nil {
        coordinate = &AdvertCoordinate{
            Longitude: longitude,
            Latitude: latitude,
        }
    } else {
        *errors = append(*errors, fmt.Errorf("ads/ad/positions/coordinate"))
    }

    if locs := ad.FindElements("./loc:locations/loc:location"); locs != nil {
        for i, loc := range locs {
            locId := FallbackUintWithReport(
                ConvString2Uint(ExtractAttrByTag(loc, "id")))(
                0, errors, fmt.Errorf("ads/ad/positions/locations[%d]/id", i))

            locName := FallbackStringWithReport(
                ExtractText(loc, "./loc:localized-name"))(
                "", errors, fmt.Errorf("ads/ad/positions/location[%d]/name", i))

            locParentId, _ := ConvString2Uint(ExtractText(loc, "./loc:parent-id"))

            locations = append(locations, AdvertLocation{
                ID:         *ReplaceUintWithNil(&locId, 0),
                Name:       *ReplaceStringWithNil(&locName, ""),
                ParentID:   ReplaceUintWithNil(&locParentId, 0),
            })
        }
    }

    return &AdvertPosition{
        Address:    ReplaceStringWithNil(&address, ""),
        City:       ReplaceStringWithNil(&city, ""),
        State:      ReplaceStringWithNil(&state, ""),
        Country:    ReplaceStringWithNil(&country, ""),
        Coordinate: coordinate,
        Locations:  locations,
    }
}

func buildPicture(ad *etree.Element, errors *[]error, _ *bool) []AdvertPicture {
    var pictures []AdvertPicture

    if pics := ad.FindElements("./pic:pictures/pic:picture"); pics != nil {
        for i, pic := range pics {
            if pic != nil {
                thumbnail := FallbackStringWithReport(
                    ExtractAttrByTag(pic.FindElement("./pic:link[@rel='thumbnail']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/thumbnail", i))
                normal := FallbackStringWithReport(
                    ExtractAttrByTag(pic.FindElement("./pic:link[@rel='normal']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/normal", i))
                large := FallbackStringWithReport(
                    ExtractAttrByTag(pic.FindElement("./pic:link[@rel='large']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/large", i))
                extraLarge := FallbackStringWithReport(
                    ExtractAttrByTag(pic.FindElement("./pic:link[@rel='extraLarge']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/extraLarge", i))
                extra2XLarge := FallbackStringWithReport(
                    ExtractAttrByTag(pic.FindElement("./pic:link[@rel='extraExtraLarge']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/extra2XLarge", i))

                pictures = append(pictures, AdvertPicture{
                    Thumbnail:    ReplaceStringWithNil(&thumbnail, ""),
                    Normal:       ReplaceStringWithNil(&normal, ""),
                    Large:        ReplaceStringWithNil(&large, ""),
                    ExtraLarge:   ReplaceStringWithNil(&extraLarge, ""),
                    Extra2XLarge: ReplaceStringWithNil(&extra2XLarge, ""),
                })
            }
        }
    }

    return pictures
}

func buildAttribute(ad *etree.Element, errors *[]error, _ *bool) []AdvertAttribute {
    var attributes []AdvertAttribute

    if attrs := ad.FindElements("./attr:attributes/attr:attribute"); attrs != nil {
        for i, attr := range attrs {
            if attr != nil {
                keySlug := FallbackStringWithReport(
                    ExtractAttrByTag(attr, "name"))(
                    "", errors, fmt.Errorf("ads/ad/attributes[%d]/key_slug", i))
                keyName := FallbackStringWithReport(
                    ExtractAttrByTag(attr, "localized-label"))(
                    "", errors, fmt.Errorf("ads/ad/attributes[%d]/key_name", i))
                valueType := FallbackStringWithReport(
                    ExtractAttrByTag(attr, "type"))(
                    "", errors, fmt.Errorf("ads/ad/attributes[%d]/value_type", i))

                valueSlugStr, _ := ExtractText(attr, "./attr:value")
                valueSlug := ReplaceStringWithNil(&valueSlugStr, "")

                valueNameStr, _ := ExtractAttrByTag(attr.FindElement("./attr:value"), "localized-label")
                valueName := ReplaceStringWithNil(&valueNameStr, "")

                attributes = append(attributes, AdvertAttribute{
                    KeySlug:    *ReplaceStringWithNil(&keySlug, ""),
                    KeyName:    *ReplaceStringWithNil(&keyName, ""),
                    ValueType:  ReplaceStringWithNil(&valueType, ""),
                    ValueSlug:  ReplaceStringWithNil(valueSlug, ""),
                    ValueName:  ReplaceStringWithNil(valueName, ""),
                })
            }
        }
    }

    return attributes
}

func buildTimestamp(ad *etree.Element, _ *[]error, _ *bool) AdvertTimestamp {
    advertCreationTime := formatTimestamp(ad.FindElement("./ad:creation-date-time"))
    advertModificationTime := formatTimestamp(ad.FindElement("./ad:modification-date-time"))
    advertStartTime := formatTimestamp(ad.FindElement("./ad:start-date-time"))
    advertEndTime := formatTimestamp(ad.FindElement("./ad:end-date-time"))

    return AdvertTimestamp{
        CreationTime:       advertCreationTime,
        ModificationTime:   advertModificationTime,
        StartTime:          advertStartTime,
        EndTime:            advertEndTime,
    }
}

func buildContact(ad *etree.Element, _ *[]error, _ *bool) *AdvertContact {
    name, _  := ExtractText(ad, "./ad:poster-contact-name")
    phone, _ := ExtractText(ad, "./ad:phone")
    phone     = strings.Replace(phone, " ", "", -1)

    return &AdvertContact{
        Name:   ReplaceStringWithNil(&name, ""),
        Phone:  ReplaceStringWithNil(&phone, ""),
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
