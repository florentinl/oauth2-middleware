package internal

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandString generates, from the set of capital and lowercase letters, a cryptographically-secure pseudo-random string of a given length.
func RandString(n int) (string, error) {
	return StringFromCharset(n, letterBytes)
}

// StringFromCharset generates, from a given charset, a cryptographically-secure pseudo-random string of a given length.
func StringFromCharset(n int, charset string) (string, error) {
	b := make([]byte, n)
	maxIdx := big.NewInt(int64(len(charset)))
	for i := 0; i < n; i++ {
		randIdx, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		// randIdx is necessarily safe to convert to int, because the max came from an int.
		randIdxInt := int(randIdx.Int64())
		b[i] = charset[randIdxInt]
	}
	return string(b), nil
}

func getBaseUri(r *http.Request) string {
	baseUri := r.Host
	if strings.Split(baseUri, ":")[0] == "localhost" {
		baseUri = "http://" + baseUri
	} else {
		baseUri = "https://" + baseUri
	}
	return baseUri
}
