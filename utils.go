package crater

import (
	"crypto/rand"
	"fmt"
)

// GenerateId generates random Id string.
func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
