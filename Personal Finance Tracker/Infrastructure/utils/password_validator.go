package utils

import (
    "errors"
    "regexp"
)

func IsValidPassword(password string) error {
    // Check minimum length
    if len(password) < 8 {
        return errors.New("password too short: must be at least 8 characters")
    }

    // Check for at least one lowercase letter
    lowercaseRegex := regexp.MustCompile(`[a-z]`)
    if !lowercaseRegex.MatchString(password) {
        return errors.New("missing lowercase letter")
    }

    // Check for at least one uppercase letter
    uppercaseRegex := regexp.MustCompile(`[A-Z]`)
    if !uppercaseRegex.MatchString(password) {
        return errors.New("missing uppercase letter")
    }

    // Check for at least one digit
    digitRegex := regexp.MustCompile(`\d`)
    if !digitRegex.MatchString(password) {
        return errors.New("missing digit")
    }

    // Check for at least one special character
    specialRegex := regexp.MustCompile(`[@$!%*?&#]`)
    if !specialRegex.MatchString(password) {
        return errors.New("missing special character (@$!%*?&#)")
    }

    // Check that password only contains allowed characters
    allowedRegex := regexp.MustCompile(`^[A-Za-z\d@$!%*?&#]+$`)
    if !allowedRegex.MatchString(password) {
        return errors.New("contains invalid characters")
    }

    return nil
}
