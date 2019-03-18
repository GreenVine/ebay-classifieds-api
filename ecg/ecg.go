package ecg

import (
    "github.com/GreenVine/ebay-ecg-api/ecg/models"
    . "github.com/GreenVine/ebay-ecg-api/utils"
    "github.com/beevik/etree"
    "github.com/parnurzeal/gorequest"
    "time"
)

// RequestEndpoint will send request to ECG API endpoint
func (agent Agent) RequestEndpoint(url string, timeout time.Duration) (*etree.Document, *ecgmodels.EndpointErrorResponse) {
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
        xml, err := ParseXML(body)
        statusCode = uint(resp.StatusCode)

        if err == nil { // XML is valid
            if root := xml.Root(); statusCode == 200 && root != nil && root.Tag != "api-base-error" && root.Tag != "html" {
                return xml, nil
            } else {
                extractedMsg, _ := ExtractText(root, "//message")
                errMsg          := ReplaceStringWithNil(&extractedMsg, "")

                return nil, &ecgmodels.EndpointErrorResponse{
                    StatusCode: &statusCode,
                    Message:    errMsg,
                }
            }
        } else {
            errMsg := "Internal server error"

            return nil, &ecgmodels.EndpointErrorResponse{ // failed to parse response
                StatusCode: &statusCode,
                Message:    &errMsg,
            }
        }
    } else {
        return nil, &ecgmodels.EndpointErrorResponse{
            StatusCode: &statusCode,
            Message:    &errMsg,
        }
    }
}
