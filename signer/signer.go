package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

type Signer struct{ Secret []byte }

func (s *Signer) Sign(method, path string, exp time.Time) string {
	u := url.Values{}
	u.Set("exp", fmt.Sprintf("%d", exp.Unix()))
	mac := hmac.New(sha256.New, s.Secret)
	mac.Write([]byte(method + "\n" + path + "\n" + u.Get("exp")))
	u.Set("sig", base64.RawURLEncoding.EncodeToString(mac.Sum(nil)))
	return u.Encode()
}
