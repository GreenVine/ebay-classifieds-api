package ecg_test

import (
    "encoding/json"
    "fmt"
    "github.com/GreenVine/ebay-classifieds-api"
    "github.com/GreenVine/ebay-classifieds-api/parsers/au"
)

var agent ecg.Agent

func ExampleAgent() {
    agent = ecg.Agent {
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
}

func ExampleAgent_RequestEndpoint() {
    advertisement, err := agent.RequestEndpoint("/ads/123456", 2000)
    category, err := agent.RequestEndpoint("/ads", 10000)

    if err != nil { // erroneous HTTP response
        fmt.Println(*err.StatusCode, *err.Message)
    } else { // successful response
        advert, errs, isFatal := auparser.ParseAdvert(advertisement) // parse an Advertisement
        cat, errs, isFatal := auparser.ParseCategory(category) // parse Advertisements in a category

        if errs != nil && !isFatal {
            fmt.Println(advert.Title) // get advertisement title
            fmt.Println(cat.Pagination.CurrentPage) // get current page number

            jsonstr, _ := json.Marshal(advert)
            fmt.Println(string(jsonstr)) // or optionally print the entire advertisement response as JSON
        }
    }
}
