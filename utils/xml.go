package utils

import (
    "fmt"
    "github.com/beevik/etree"
)

// ParseXML will parse XML responses and convert it to JSON
func ParseXML(rawXML string) (*etree.Document, error) {
    doc := etree.NewDocument()

    if err := doc.ReadFromString(rawXML); err != nil {
        return nil, fmt.Errorf("invalid or malformed XML")
    }

    return doc, nil
}

// ExtractElementText is a safe wrapper of FindElement with Text() applied on the result
func ExtractElementText(element *etree.Element, path string) (string, error) {
    if result := element.FindElement(path); result != nil {
        return result.Text(), nil
    }

    return "", fmt.Errorf("element not found")
}

// ExtractElementsText is a safe wrapper of FindElements with Text() applied on results
func ExtractElementsText(element *etree.Element, path string) ([]string, error) {
    var texts []string

    if results := element.FindElements(path); results != nil {
        for _, result := range results {
            if result != nil {
                texts = append(texts, result.Text())
            }
        }

        return texts, nil
    }

    return texts, fmt.Errorf("element not found")
}
