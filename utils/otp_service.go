package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateOTP generates a 6-digit OTP using cryptographically secure random numbers
func GenerateOTP() (int, error) {
	var otp int
	for i := 0; i < 6; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return 0, err
		}
		otp = otp*10 + int(digit.Int64())
	}
	return otp, nil
}
