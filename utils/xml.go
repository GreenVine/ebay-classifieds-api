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

// ExtractText is a safe wrapper of FindElement with Text() applied on the result
func ExtractText(element *etree.Element, path string) (string, error) {
    if element != nil {
        if result := element.FindElement(path); result != nil {
            return result.Text(), nil
        }
    }

    return "", fmt.Errorf("element invalid or does not exist")
}

// ExtractTextAsUint wraps ExtractText and converts result to uint
func ExtractTextAsUint(element *etree.Element, path string) (uint, error) {
    return ConvString2Uint(ExtractText(element, path))
}

// ExtractTextAsFloat64 wraps ExtractText and converts result to uint
func ExtractTextAsFloat64(element *etree.Element, path string) (float64, error) {
    return ConvString2Float64(ExtractText(element, path))
}

// ExtractAttrByTag extracts a given tag from attributes
func ExtractAttrByTag(element *etree.Element, tag string) (string, error) {
    if element != nil {
        if result := element.SelectAttr(tag); result != nil {
            return result.Value, nil
        }
    }

    return "", fmt.Errorf("attribute invalid or does not exist")
}
