package utils

import "encoding/hex"

// Contains is a helper function to check if string is in a string slice
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// IsCommitHash returns whether the given string is a valid sha-1 commit ID,
// this means 40 chars long, with only hexadecimal digits
func IsCommitHash(s string) bool {
	if len(s) != 40 {
		return false
	}

	// try to convert it from hexadecimal string to a byte array
	_, err := hex.DecodeString(s)
	return err == nil
}
