package ecg

import (
    "github.com/GreenVine/ebay-ecg-api/ecg"
    u "github.com/GreenVine/ebay-ecg-api/utils"
    "github.com/beevik/etree"
    "github.com/parnurzeal/gorequest"
    "time"
)

// Agent is ECG agent that stores configurable information
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

func (agent Agent) hasECGAuthorization() bool {
    return agent.ECGAuthorization != nil
}

func (agent Agent) hasECGAuthentication() bool {
    return agent.ECGAuthentication != nil
}

// RequestEndpoint will send request to ECG API endpoint
func (agent Agent) RequestEndpoint(url string, timeout time.Duration) (*etree.Document, *ecg.EndpointErrorResponse) {
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

                return nil, &ecg.EndpointErrorResponse{
                    StatusCode: &statusCode,
                    Message:    errMsg,
                }
            }

            return xml, nil
        }

        errMsg := "Internal server error"

        return nil, &ecg.EndpointErrorResponse{ // failed to parse response
            StatusCode: &statusCode,
            Message:    &errMsg,
        }
    }

    return nil, &ecg.EndpointErrorResponse{
        StatusCode: &statusCode,
        Message:    &errMsg,
    }
}
