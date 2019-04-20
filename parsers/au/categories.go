package auparser

import (
    "fmt"
    models "github.com/GreenVine/ebay-classifieds-api/parsers/au/models"
    u "github.com/GreenVine/ebay-classifieds-api/utils"
    "github.com/beevik/etree"
)

// ParseCategories is to build a Categories model from raw XML response
func ParseCategories(doc *etree.Document) (*models.Categories, []error, bool) {
    if doc == nil {
        return nil, []error{ fmt.Errorf("empty API response") }, true
    }

    root := doc.Root()

    if root == nil || root.Space != "cat" || root.Tag != "categories" {
       return nil, []error{ fmt.Errorf("unexpected API response") }, true
    }

    var errors []error
    var hasCriticalError = false

    if rootCategory := root.SelectElement("category"); rootCategory != nil {
        if categories := buildCategories(rootCategory, &errors, &hasCriticalError); !hasCriticalError && categories != nil {
            return categories, errors, false
        }
    }

    return nil, errors, true
}

func buildCategories(category *etree.Element, errors *[]error, hasCriticalError *bool) *models.Categories {
    if category != nil {
        catID, err := u.ConvString2Uint(u.ExtractAttrByTag(category, "id"))
        if err != nil {
           return nil
        }

        catName := u.FallbackStringWithReport(
           u.ExtractText(category, "./cat:localized-name"))(
           "", errors, fmt.Errorf("categories/category/%d/name", catID))

        catSlug := u.FallbackStringWithReport(
           u.ExtractText(category, "./cat:id-name"))(
           "", errors, fmt.Errorf("categories/category/%d/slug", catID))

        catParentID := u.FallbackUintWithReport(
            u.ExtractTextAsUint(category, "./cat:parent-id"))(
            0, errors, fmt.Errorf("categories/category/%d/parent_id", catID))

        catParentSlug := u.FallbackStringWithReport(
           u.ExtractText(category, "./cat:l1-name"))(
           "", errors, fmt.Errorf("categories/category/%d/parent_slug", catID))

        catChildrenCount := u.FallbackUintWithReport(
           u.ExtractTextAsUint(category, "./cat:children-count"))(
           0, errors, fmt.Errorf("categories/category/%d/children_count", catID))

        var subcategories []models.Categories

        if subcategoriesList := category.FindElements("./cat:category"); subcategoriesList != nil {
            // recursively add subcategories

            for _, subcategory := range subcategoriesList {
                subcategories = append(subcategories, *buildCategories(subcategory, errors, hasCriticalError))
            }
        }

        return &models.Categories{
            ID:             catID,
            Name:           catName,
            Slug:           catSlug,
            ParentID:       &catParentID,
            ParentSlug:     &catParentSlug,
            ChildrenCount:  catChildrenCount,
            Subcategories:  subcategories,
            IsRootCategory: catID <= 0,
        }
    }

    return nil
}
