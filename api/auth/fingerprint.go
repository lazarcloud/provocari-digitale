package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strings"
)

func GetIPAddress(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-IP")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	ipAddress = strings.Split(ipAddress, ":")[0]

	return ipAddress
}

func GenerateFingerprint(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	ipAddress := GetIPAddress(r)
	sessionID, _ := r.Cookie("session_id") // Get session ID from cookies

	combinedInfo := userAgent + ipAddress + sessionID.Value

	hasher := sha1.New()
	hasher.Write([]byte(combinedInfo))
	fingerprint := hex.EncodeToString(hasher.Sum(nil))

	return fingerprint
}
