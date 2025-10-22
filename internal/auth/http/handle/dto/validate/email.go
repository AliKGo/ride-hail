package validate

import (
	"net"
	"net/mail"
	"strings"
)

func isValidEmailBasic(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" || len(email) > 254 {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

// HasMX checks if domain has MX or A/AAAA records (fallback)
func hasMX(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	// MX check
	mxRecords, err := net.LookupMX(domain)
	if err == nil && len(mxRecords) > 0 {
		return true
	}

	// Fallback to A/AAAA
	_, errA := net.LookupHost(domain)
	return errA == nil
}

func ValidateEmail(email string, requireMX bool) bool {
	if !isValidEmailBasic(email) {
		return false
	}
	if requireMX {
		// делает сетевой вызов — учитывать время ожидания/контекст в реальном коде
		return hasMX(email)
	}
	return true
}
