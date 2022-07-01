// Generators of fake byte arrays for test purpose
package fakeData

import (
	cryptoRand "crypto/rand"
	"fmt"
)

// Generate array of random byte arrays for testing
func GetByteArrays(numOfData int, dataLength int) ([][]byte, error) {
	var testStrings [][]byte

	for i := 0; i < numOfData; i++ {
		testString, err := GetByteArray(dataLength)
		if err != nil {
			return nil, err
		}

		testStrings = append(testStrings, testString)
	}

	return testStrings, nil
}

// Generate random data array for testing
func GetByteArray(lengthOfData int) ([]byte, error) {
	message := make([]byte, lengthOfData)

	if _, err := cryptoRand.Read(message); err != nil {
		return nil, fmt.Errorf("GetByteArray error: %w", err)
	}

	return message, nil
}
