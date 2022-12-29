package fakeData

import (
	"fmt"
	"math/rand"
	"strings"
)

// Get random test string with length
// Stole from https://github.com/jaswdr/faker/blob/master/faker.go
func GetRandomStringWithLength(lengthOfString int) string {
	r := []string{}
	for i := 0; i < lengthOfString; i++ {
		r = append(r, GetRandomLetter())
	}

	return strings.Join(r, "")
}

// nolint
func GetRandomLetter() string {
	cNum := 97 + rand.Intn(25)

	return fmt.Sprintf("%c", cNum)
}
