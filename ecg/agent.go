package ecg

// Agent is ECG agent that stores configurable information
type Agent struct {
    Endpoint string
    ECGAuthorization *Authorization
    ECGAuthentication *Authentication
}

type Authentication struct {
    AuthenticateUser string
    AuthenticateAd string
    AuthenticateDevice string
}

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
