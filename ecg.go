package ecg

import (
    u "github.com/GreenVine/ebay-classifieds-api/utils"
    "github.com/beevik/etree"
    "github.com/parnurzeal/gorequest"
    "time"
)

// Agent initialises the ECG Agent and pass the API endpoint as base URL.
// Optionally, you can pass either authentication or authorization or both settings to the agent.
//
// Read more about security settings implemented in the API here: https://api.ebay-kleinanzeigen.de/docs/pages/security.
type Agent struct {
    Endpoint string
    ECGAuthorization *Authorization
    ECGAuthentication *Authentication
}

// Authentication is ECG authentication settings
type Authentication struct {
    AuthenticateUser string
    AuthenticateAd string
    AuthenticateDevice string
}

// Authorization is ECG authorization settings
type Authorization struct {
    Username string
    Password string
}

// EndpointErrorResponse is model of erroneous endpoint response
type EndpointErrorResponse struct {
    StatusCode  *uint   `json:"code"`
    Message     *string `json:"message"`
}

func (agent Agent) hasECGAuthorization() bool {
    return agent.ECGAuthorization != nil
}

func (agent Agent) hasECGAuthentication() bool {
    return agent.ECGAuthentication != nil
}

// RequestEndpoint requests the API endpoint along with the URL and timeout (in milliseconds) settings
// ECG Agent will either return a XML document on success, or an `EndpointErrorResponse` type on failure.
func (agent Agent) RequestEndpoint(url string, timeout time.Duration) (*etree.Document, *EndpointErrorResponse) {
    http := gorequest.
        New().
        Timeout(timeout * time.Millisecond)

    if agent.hasECGAuthorization() {
        http.SetBasicAuth((*agent.ECGAuthorization).Username, (*agent.ECGAuthorization).Password)
    }

    resp, body, errs := http.Get(agent.Endpoint + url).End()

    var statusCode uint = 503 // error by default
    errMsg := "Service temporarily unavailable"

    if errs == nil && len(errs) < 1 && resp != nil && body != "" {
        xml, err := u.ParseXML(body)
        statusCode = uint(resp.StatusCode)

        if err == nil { // XML is valid
            if root := xml.Root(); statusCode != 200 || root == nil || root.Tag == "api-base-error" || root.Tag == "html" {
                extractedMsg, _ := u.ExtractText(root, "//message")
                errMsg          := u.ReplaceStringWithNil(&extractedMsg, "")

                return nil, &EndpointErrorResponse{
                    StatusCode: &statusCode,
                    Message:    errMsg,
                }
            }

            return xml, nil
        }

        errMsg := "Internal server error"

        return nil, &EndpointErrorResponse{ // failed to parse response
            StatusCode: &statusCode,
            Message:    &errMsg,
        }
    }

    return nil, &EndpointErrorResponse{
        StatusCode: &statusCode,
        Message:    &errMsg,
    }
}
