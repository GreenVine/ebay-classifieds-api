package ecg

// EndpointErrorResponse is model of erroneous endpoint response
type EndpointErrorResponse struct {
    StatusCode  *uint   `json:"code"`
    Message     *string `json:"message"`
}
