package auparsers

import (
    "fmt"
    . "github.com/GreenVine/ebay-ecg-api/processors/au/models"
    . "github.com/GreenVine/ebay-ecg-api/utils"
    "github.com/beevik/etree"
)

// ParseCategory is to build a Category models from raw XML response
func ParseCategory(doc *etree.Document) (*Category, []error, bool) {
    root := doc.Root()

    if root == nil || root.Space != "ad" || root.Tag != "ads" {
       return nil, []error{ fmt.Errorf("unexpected API response") }, true
    }

    var errors []error
    var hasCriticalError = false

    if category := buildCategoryBase(root, &errors, &hasCriticalError); hasCriticalError {
        return nil, errors, true
    } else {
        return &category, errors, false
    }
}

func buildCategoryBase(root *etree.Element, errors *[]error, hasCriticalError *bool) Category {
    var adverts []Advert

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

    return Category{
        Adverts:        adverts,
        Pagination:     pagination,
    }
}

func buildPagination(root *etree.Element, errors *[]error, _ *bool) *CategoryPagination {
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

      return &CategoryPagination{
          CurrentPage: currentPage,
          PageSize:    pageSize,
          EntrySize:   matchedEntries,
      }
    }

    return nil
}
