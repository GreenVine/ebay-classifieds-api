package ecg

import (
    "fmt"
    "github.com/parnurzeal/gorequest"
    "time"
)

// RequestAPI will send request to ECG API endpoint
func (agent Agent) RequestAPI(url string, timeout time.Duration) (string, []error) {
    http := gorequest.
        New().
        Timeout(timeout * time.Millisecond)

    if agent.hasECGAuthorization() {
        http.SetBasicAuth((*agent.ECGAuthorization).Username, (*agent.ECGAuthorization).Password)
    }

    resp, body, errs := http.Get(agent.Endpoint + url).End()

    if errs != nil || resp.StatusCode != 200 {
        errs = append(errs, fmt.Errorf("unexpected HTTP status code %d received", resp.StatusCode))
        return "", errs
    }

    return body, nil
}
