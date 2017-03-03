package sendowl

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
)

const SignatureQueryParam = "signature"

func Signature(signingKey, signingKeySecret string, vs url.Values) string {
	return base64.StdEncoding.EncodeToString(generateSignature(signingKey, signingKeySecret, vs))
}

func generateSignature(signingKey, signingKeySecret string, vs url.Values) []byte {
	mac := hmac.New(sha1.New, []byte(SigningKey(signingKey, signingKeySecret)))
	mac.Write([]byte(SigningText(signingKeySecret, vs)))
	return mac.Sum(nil)
}

func VerifySignatureFromRequest(signingKey, signingKeySecret string, r *http.Request) bool {
	expectedMAC := []byte(Signature(signingKey, signingKeySecret, r.URL.Query()))
	givenMAC := []byte(r.URL.Query().Get(SignatureQueryParam))
	return hmac.Equal(expectedMAC, givenMAC)
}

func SigningText(signingKeySecret string, vs url.Values) string {
	vs.Del("signature")
	return fmt.Sprintf("%s&secret=%s", vs.Encode(), url.QueryEscape(signingKeySecret))
}

func SigningKey(signingKey, signingKeySecret string) string {
	return fmt.Sprintf("%s&%s", signingKey, signingKeySecret)
}
