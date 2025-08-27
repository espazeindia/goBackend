package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

// GenerateOTP generates a 6-digit OTP using cryptographically secure random numbers
func GenerateOTP() (int, error) {
	var otp int
	for len(strconv.Itoa(otp)) < 6 {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return 0, err
		}
		otp = otp*10 + int(digit.Int64())
	}
	return otp, nil
}
