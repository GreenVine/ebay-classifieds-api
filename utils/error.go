package utils

// FallbackUintWithReport is to replace the variable with the fallback value if an error exists
// and append the error to the array
func FallbackUintWithReport(value uint, err error) func(fallback uint, errArray *[]error, displayError error) uint {
    return func(fallback uint, errArray *[]error, displayError error) uint {
        if err != nil { // fallback is necessary
            *errArray = append(*errArray, displayError)
            return fallback
        }

        return value // return actual value as no error occurs
    }
}

// FallbackFloat64WithReport is to replace the variable with the fallback value if an error exists
// and append the error to the array
func FallbackFloat64WithReport(value float64, err error) func(fallback float64, errArray *[]error, displayError error) float64 {
    return func(fallback float64, errArray *[]error, displayError error) float64 {
        if err != nil { // fallback is necessary
            *errArray = append(*errArray, displayError)
            return fallback
        }

        return value // return actual value as no error occurs
    }
}

// FallbackStringWithReport is to replace the variable with the fallback value if an error exists
// and append the error to the array
func FallbackStringWithReport(value string, err error) func(fallback string, errArray *[]error, displayError error) string {
    return func(fallback string, errArray *[]error, displayError error) string {
        if err != nil { // fallback is necessary
            *errArray = append(*errArray, displayError)
            return fallback
        }

        return value // return actual value as no error occurs
    }
}

// ReplaceStringWithNil replaces a matched string with nil
func ReplaceStringWithNil(str *string, match string) *string {
    if str != nil && *str != match {
        return str
    }

    var fallbackVal string
    return &fallbackVal
}

// ReplaceUintWithNil replaces a matched unit with nil
func ReplaceUintWithNil(number *uint, match uint) *uint {
    if number != nil && *number != match {
        return number
    }

    var fallbackVal uint
    return &fallbackVal
}
