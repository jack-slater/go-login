package helpers

import (
	"crypto/sha256"
	"fmt"
)

func IncryptPassword(p string) string  {
	sum := sha256.Sum256([]byte(p))
	return fmt.Sprintf("%x", sum)
}
