package storage

import (
	"codebase-app/internal/infrastructure/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func GenerateSignedURL(filename string, expiration time.Duration) string {
	urlToSigned := config.Envs.App.BaseURL + "/api/storage/private/" + filename
	var (
		key            = []byte(config.Envs.Guard.JwtPrivateKey)
		expirationTime = time.Now().UTC().Add(expiration).Unix()
		data           = fmt.Sprintf("%s%d", urlToSigned, expirationTime)
	)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))

	// Add the expiration time and signature to the URL
	u, _ := url.Parse(urlToSigned)
	q := u.Query()
	q.Set("expires", strconv.FormatInt(expirationTime, 10))
	q.Set("signature", signature)
	u.RawQuery = q.Encode()

	return u.String()
}
