package utils

import "strconv"

// ConvString2Uint is a wrapper to safely convert string to uint
func ConvString2Uint(text string, err1 error) (uint, error) {
    if err1 == nil {
        if parsedUint, err2 := strconv.ParseUint(text, 10, 64); err2 != nil {
            return 0, err2
        } else {
            return uint(parsedUint), nil
        }
    } else {
        return 0, err1
    }
}

// ConvString2Float64 is a wrapper to safely convert string to uint
func ConvString2Float64(text string, err1 error) (float64, error) {
    if err1 == nil {
        if parsedFloat64, err2 := strconv.ParseFloat(text, 64); err2 != nil {
            return 0, err2
        } else {
            return parsedFloat64, nil
        }
    } else {
        return 0, err1
    }
}
