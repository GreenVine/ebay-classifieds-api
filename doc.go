// Package ecg is an unofficial SDK for eBay Classifieds (ECG) API written in Go. To use the API, you must have partner access.
// 
// The module primarily consists of two components, ECG Agent and ECG Parser. The former sends the request to the API endpoint, whereas the latter parses the platform-dependent responses received from the endpoint. Currently, it implements APIs available in Australia but may also suitable for use in other countries.
// 
// Official API documentation provided by eBay Germany can be found here: https://api.ebay-kleinanzeigen.de/docs/pages/home.
// 
// Then import it within your project files whenever you need ECG Agent and/or ECG Parser (country-specific):
// 
//     import (
//         "github.com/GreenVine/ebay-classifieds-api"
//         "github.com/GreenVine/ebay-classifieds-api/parsers/au"
//     )
package ecg
