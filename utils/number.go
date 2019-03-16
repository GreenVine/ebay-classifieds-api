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

// ConvString2UintDefault is a wrapper to ConvString2Uint that returns fallback value if it fails
func ConvString2UintDefault(text string, err1 error) func(fallback uint) uint {
    return func(fallback uint) uint {
        if parsedUint, err2 := ConvString2Uint(text, err1); err2 != nil {
            return fallback
        } else {
            return parsedUint
        }
    }
}
