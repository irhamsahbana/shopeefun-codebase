package pkg

import (
	"crypto/rand"
	"math/big"
)

// generates a random string of specified length
func GeneratePassword(length int) string {
	const (
		uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
		specialChars   = "!@#$%^&*()_+{}[]:;<>,.?/~"
		numericChars   = "0123456789"
	)

	var (
		uppercaseIdx, lowercaseIdx, specialIdx, numericIdx int
		randomIdx                                          int
		result                                             string
	)

	// Ensure at least one of each type of character
	result += string(uppercaseChars[randInt(0, len(uppercaseChars))])
	result += string(lowercaseChars[randInt(0, len(lowercaseChars))])
	result += string(specialChars[randInt(0, len(specialChars))])
	result += string(numericChars[randInt(0, len(numericChars))])

	for i := 4; i < length; i++ {
		randomIdx = randInt(0, 4)
		switch randomIdx {
		case 0:
			uppercaseIdx = randInt(0, len(uppercaseChars))
			result += string(uppercaseChars[uppercaseIdx])
		case 1:
			lowercaseIdx = randInt(0, len(lowercaseChars))
			result += string(lowercaseChars[lowercaseIdx])
		case 2:
			specialIdx = randInt(0, len(specialChars))
			result += string(specialChars[specialIdx])
		case 3:
			numericIdx = randInt(0, len(numericChars))
			result += string(numericChars[numericIdx])
		}
	}
	return result
}

// randInt returns a random integer in the range [min, max)
func randInt(min, max int) int {
	if min == max {
		return min
	}
	maxBigInt := big.NewInt(int64(max) - int64(min))
	randBigInt, _ := rand.Int(rand.Reader, maxBigInt)
	return min + int(randBigInt.Int64())
}
