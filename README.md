# ebay-classifieds-api

This is an unofficial SDK for eBay Classifieds (ECG) API written in Go. To use the API, you must have partner access.

The module primarily consists of two components, ECG Agent and ECG Parser. The former sends the request to the API endpoint, whereas the latter parses the platform-dependent responses received from the endpoint. Currently, it implements APIs available in Australia but may also suitable for use in other countries.

Official API documentation provided by eBay Germany can be found [here](https://api.ebay-kleinanzeigen.de/docs/pages/home).

## Installation

To install the module, simply run:

```bash
go get github.com/GreenVine/ebay-classifieds-api
```

Then import it within your project files whenever you need ECG Agent and/or ECG Parser (country-specific):

```go
import (
    "github.com/GreenVine/ebay-classifieds-api"
    "github.com/GreenVine/ebay-classifieds-api/parsers/au"
)
```

## Usage

Initialise the ECG Agent as follows and pass the API endpoint as base URL. Optionally, you can pass either authentication or authorization or both settings to the agent. Read more about security settings implemented in the API [here](https://api.ebay-kleinanzeigen.de/docs/pages/security).

```go
ecg := ecg.Agent {
    Endpoint: "https://api.example.com/api", // API Endpoint
    ECGAuthentication: &ecg.Authentication{ // Authentication (optional)
        AuthenticateUser: "user",
        AuthenticateAd: "ad",
        AuthenticateDevice: "device",
    },
    ECGAuthorization: &ecg.Authorization{ // Authorization via HTTP Header (optional)
        Username: "user",
        Password: "password",
    },
}
```

Then request the API along with the URL and timeout (in milliseconds) settings:

```go
advertisement, err := ecg.RequestEndpoint("/ads/123456", 2000)
category, err := ecg.RequestEndpoint("/ads", 10000)
```

ECG Agent will either return a XML document on success, or an `EndpointErrorResponse` type on failure. You will then need to use a country-specific parser to parse the advertisement or category response:

```go
if err != nil { // erroneous HTTP response
   fmt.Println(*err.StatusCode, *err.Message)
} else { // successful response
   advert, errs, isFatal := auparser.ParseAdvert(advertisement) // parse an Advertisement
   cat, errs, isFatal := auparser.ParseCategory(category) // parse Advertisements in a category

   fmt.Println(advert.Title) // get advertisement title
   fmt.Println(cat.Pagination.CurrentPage) // get current page number
   
   jsonstr, _ := json.Marshal(advert)
   fmt.Println(string(jsonstr)) // or optionally print the entire advertisement response as JSON
}
```
