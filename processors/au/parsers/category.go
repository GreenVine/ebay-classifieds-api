package auparsers

import (
    "fmt"
    models "github.com/GreenVine/ebay-ecg-api/processors/au/models"
    . "github.com/GreenVine/ebay-ecg-api/utils"
    "github.com/beevik/etree"
    "time"
)

// ParseCategory is to build a Category models from raw XML response
func ParseCategory(rawXML string) (*models.Category, []error, bool) {
    doc, err := ParseXML(rawXML)

    if err != nil {
        return nil, []error{ err }, true
    }

    root := doc.Root()

    if root.Space != "ad" || root.Tag != "ads" {
        return nil, []error{ fmt.Errorf("unexpected API response") }, true
    }

    var errors []error
    var hasCriticalError = false

    if category := buildCategory(root, &errors, &hasCriticalError); hasCriticalError {
        return nil, errors, true
    } else {
        return &category, errors, false
    }
}

func buildCategory(root *etree.Element, errors *[]error, hasCriticalError *bool) models.Category {
    // build each advertisement
    adverts := buildAdvert(root.SelectElements("ad"), errors, hasCriticalError)

    // build pagination
    pagination := buildPagination(root, errors, hasCriticalError)

    return models.Category{
        Adverts:        adverts,
        Pagination:     pagination,
    }
}

func buildAdvert(ads []*etree.Element, errors *[]error, hasCriticalError *bool) []models.NormalisedAdvert {
    var adverts []models.NormalisedAdvert

    if ads != nil {
        for _, ad := range ads {
            if ad != nil {
                advertId, err := ConvString2Uint(ExtractAttrByTag(ad, "id"))
                if err != nil {
                    *errors = append(*errors, fmt.Errorf("ads/ad/id"))
                    *hasCriticalError = true
                    return adverts
                }

                advertType := FallbackStringWithReport(
                    ExtractText(ad, "./ad:ad-type/ad:value"))(
                    "UNKNOWN", errors, fmt.Errorf("ads/ad/type"))

                advertPrice := buildAdvertPrice(ad, errors, hasCriticalError)

                advertStatus := FallbackStringWithReport(
                    ExtractText(ad, "./ad:ad-status/ad:value"))(
                    "UNKNOWN", errors, fmt.Errorf("ads/ad/status"))

                advertCategory := buildAdvertCategory(ad, errors, hasCriticalError)

                advertPosition := buildPosition(ad, errors, hasCriticalError)

                advertPosterType := FallbackStringWithReport(
                    ExtractText(ad, "./ad:poster-type/ad:value"))(
                    "", errors, fmt.Errorf("ads/ad/poster_type"))

                advertTitle := FallbackStringWithReport(
                    ExtractText(ad, "./ad:title"))(
                    "", errors, fmt.Errorf("ads/ad/title"))

                advertDescriptionExcerptHTML := FallbackStringWithReport(
                    ExtractText(ad, "./ad:description"))(
                    "", errors, fmt.Errorf("ads/ad/desc_excerpt_html"))

                advertDescriptionExcerpt, err := FormatHtml2Base64(advertDescriptionExcerptHTML)
                if err != nil {
                    *errors = append(*errors, fmt.Errorf("ads/ad/desc_excerpt_plain_b64"))
                }

                advertPictures := buildPicture(ad, errors, hasCriticalError)

                advertAttributes := buildAttribute(ad, errors, hasCriticalError)

                advertCreationTime := buildTime(ad.FindElement("./ad:creation-date-time"), errors, hasCriticalError)

                advertStartTime := buildTime(ad.FindElement("./ad:start-date-time"), errors, hasCriticalError)

                advertEndTime := buildTime(ad.FindElement("./ad:end-date-time"), errors, hasCriticalError)

                adverts = append(adverts, models.NormalisedAdvert{
                    ID:                     advertId,
                    Type:                   advertType,
                    Status:                 advertStatus,
                    Category:               advertCategory,
                    Position:               advertPosition,
                    PosterType:             ReplaceStringWithNil(advertPosterType, ""),
                    Price:                  advertPrice,
                    Title:                  *ReplaceStringWithNil(advertTitle, ""),
                    DescriptionExcerptB64:  ReplaceStringWithNil(advertDescriptionExcerpt, ""),
                    DescriptionExcerptHTML: ReplaceStringWithNil(advertDescriptionExcerptHTML, ""),
                    Pictures:               advertPictures,
                    Attributes:             advertAttributes,
                    CreationTime:           advertCreationTime,
                    StartTime:              advertStartTime,
                    EndTime:                advertEndTime,
                })
            }

        }
    }

    return adverts
}

func buildAdvertCategory(ad *etree.Element, errors *[]error, _ *bool) *models.AdvertCategory {
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

    return &models.AdvertCategory{
        ID:             catId,
        Name:           *ReplaceStringWithNil(catName, ""),
        Slug:           ReplaceStringWithNil(catSlug, ""),
        ParentSlug:     ReplaceStringWithNil(catParentSlug, ""),
        ChildrenCount:  &catChildrenCount,
    }
}

func buildAdvertPrice(ad *etree.Element, errors *[]error, _ *bool) *models.Price {
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

    return &models.Price{
        Type: &priceType,
        Amount: &priceAmount,
        HighestAmount: &priceHighestAmount,
        Currency: ReplaceStringWithNil(currency, ""),
        CurrencySymbol: ReplaceStringWithNil(currencySymbol, ""),
    }
}

func buildPagination(root *etree.Element, errors *[]error, _ *bool) *models.Pagination {
    if root != nil {
      currentPage := FallbackUintWithReport(
          ExtractTextAsUint(root, "./ad:ads-search-options/ad:page"))(
          0, errors, fmt.Errorf("category/root/current"))
      pageSize :=  FallbackUintWithReport(
          ExtractTextAsUint(root, "./ad:ads-search-options/ad:size"))(
          0, errors, fmt.Errorf("category/root/size"))
        // retrieve total matched entries
      matchedEntries := FallbackUintWithReport(
          ExtractTextAsUint(root, "./types:paging/types:numFound"))(
            0, errors, fmt.Errorf("category/matched_entries"))

      return &models.Pagination{
          CurrentPage:      currentPage,
          PageSize:         pageSize,
          MatchedEntries:   matchedEntries,
      }
    }

    return nil
}

func buildPosition(ad *etree.Element, errors *[]error, _ *bool) *models.Position {
    var coordinate *models.Coordinate
    longitude, longerr := ConvString2Float64(ExtractText(ad, "./ad:ad-address/types:longitude"))
    latitude, laterr := ConvString2Float64(ExtractText(ad, "./ad:ad-address/types:latitude"))

    if longerr == nil && laterr == nil {
        coordinate = &models.Coordinate{
            Longitude: longitude,
            Latitude: latitude,
        }
    } else {
        *errors = append(*errors, fmt.Errorf("ads/ad/positions/coordinate"))
    }

    state := FallbackStringWithReport(
        ExtractText(ad, "./ad:ad-address/types:state"))(
        "", errors, fmt.Errorf("ads/ad/positions/state"))

    return &models.Position{
        Coordinate: coordinate,
        State: ReplaceStringWithNil(state, ""),
    }
}

func buildPicture(ad *etree.Element, errors *[]error, _ *bool) []models.Picture {
    var pictures []models.Picture

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
                extraExtraLarge := FallbackStringWithReport(
                    ExtractAttrByTag(pic.FindElement("./pic:link[@rel='extraExtraLarge']"), "href"))(
                    "", errors, fmt.Errorf("ads/ad/pictures[%d]/extraExtraLarge", i))

                pictures = append(pictures, models.Picture{
                   Thumbnail:       ReplaceStringWithNil(thumbnail, ""),
                   Normal:          ReplaceStringWithNil(normal, ""),
                   Large:           ReplaceStringWithNil(large, ""),
                   ExtraLarge:      ReplaceStringWithNil(extraLarge, ""),
                   ExtraExtraLarge: ReplaceStringWithNil(extraExtraLarge, ""),
                })
            }
        }
    }

    return pictures
}

func buildAttribute(ad *etree.Element, errors *[]error, _ *bool) []models.Attribute {
    var attributes []models.Attribute

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
                valueSlug := FallbackStringWithReport(
                    ExtractText(attr, "./attr:value"))(
                    "", errors, fmt.Errorf("ads/ad/attributes[%d]/value_slug", i))
                valueName := FallbackStringWithReport(
                   ExtractAttrByTag(attr.FindElement("./attr:value"), "localized-label"))(
                   "", errors, fmt.Errorf("ads/ad/attributes[%d]/value_name", i))

                attributes = append(attributes, models.Attribute{
                    KeySlug:    *ReplaceStringWithNil(keySlug, ""),
                    KeyName:    *ReplaceStringWithNil(keyName, ""),
                    ValueType:  ReplaceStringWithNil(valueType, ""),
                    ValueSlug:  ReplaceStringWithNil(valueSlug, ""),
                    ValueName:  ReplaceStringWithNil(valueName, ""),
                })
            }
        }
    }

    return attributes
}

func buildTime(element *etree.Element, _ *[]error, _ *bool) *time.Time {
    if element != nil {
        timestr := element.Text()

        if timeinst, err := time.Parse(time.RFC3339, timestr); err == nil {
            timeinst = timeinst.UTC()
            return &timeinst
        }
    }

    return nil
}
