package auparser

import (
    "fmt"
    models "github.com/GreenVine/ebay-classifieds-api/parsers/au/models"
    u "github.com/GreenVine/ebay-classifieds-api/utils"
    "github.com/beevik/etree"
)

// ParseCategory is to build a Category models from raw XML response
func ParseCategory(doc *etree.Document) (*models.Category, []error, bool) {
    if doc == nil {
        return nil, []error{ fmt.Errorf("empty API response") }, true
    }

    root := doc.Root()

    if root == nil || root.Space != "ad" || root.Tag != "ads" {
       return nil, []error{ fmt.Errorf("unexpected API response") }, true
    }

    var errors []error
    var hasCriticalError = false

    if category := buildCategoryBase(root, &errors, &hasCriticalError); !hasCriticalError {
        return &category, errors, false
    }

    return nil, errors, true
}

func buildCategoryBase(root *etree.Element, errors *[]error, hasCriticalError *bool) models.Category {
    var adverts []models.Advert

    for _, advert := range root.SelectElements("ad") { // build each advertisement
        if advert != nil {
            builtAdvert := BuildAdvertBase(advert, errors, hasCriticalError)

            if *hasCriticalError { // critical error that ends the entire response
                break
            } else if builtAdvert == nil { // error that skips the current ad
                continue
            } else {
                adverts = append(adverts, *builtAdvert)
            }
        }
    }

    // build pagination
    pagination := buildPagination(root, errors, hasCriticalError)

    return models.Category{
        Adverts:        adverts,
        Pagination:     pagination,
    }
}

func buildPagination(root *etree.Element, errors *[]error, _ *bool) *models.CategoryPagination {
    if root != nil {
      currentPage := u.FallbackUintWithReport(
          u.ExtractTextAsUint(root, "./ad:ads-search-options/ad:page"))(
          0, errors, fmt.Errorf("category/root/current"))
      pageSize :=  u.FallbackUintWithReport(
          u.ExtractTextAsUint(root, "./ad:ads-search-options/ad:size"))(
          0, errors, fmt.Errorf("category/root/size"))
        // retrieve total matched entries
      matchedEntries := u.FallbackUintWithReport(
          u.ExtractTextAsUint(root, "./types:paging/types:numFound"))(
            0, errors, fmt.Errorf("category/matched_entries"))

      return &models.CategoryPagination{
          CurrentPage: currentPage,
          PageSize:    pageSize,
          EntrySize:   matchedEntries,
      }
    }

    return nil
}
