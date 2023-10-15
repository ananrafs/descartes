package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func CreateHash(params ...interface{}) string {
	h := sha256.New()
	for _, param := range params {
		h.Write([]byte(fmt.Sprintf("%v", param)))
	}
	return hex.EncodeToString(h.Sum(nil))
}
