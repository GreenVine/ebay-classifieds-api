package ecgmodels

type EndpointErrorResponse struct {
    StatusCode  *uint   `json:"code"`
    Message     *string `json:"message"`
}
