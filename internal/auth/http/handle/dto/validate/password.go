package validate

import (
	"strings"
	"unicode"
)

type PasswordRules struct {
	MinLen               int
	RequireUpper         bool
	RequireLower         bool
	RequireNumber        bool
	RequireSpecial       bool
	DisallowEmailContain bool
}

var DefaultStrongPolicy = PasswordRules{
	MinLen:               12,
	RequireUpper:         true,
	RequireLower:         true,
	RequireNumber:        true,
	RequireSpecial:       true,
	DisallowEmailContain: true,
}

func isSequential(s string, length int) bool {
	if len(s) < length {
		return false
	}
	for i := 0; i <= len(s)-length; i++ {
		inc := true
		dec := true
		for j := 1; j < length; j++ {
			if s[i+j]-s[i+j-1] != 1 {
				inc = false
			}
			if s[i+j-1]-s[i+j] != 1 {
				dec = false
			}
		}
		if inc || dec {
			return true
		}
	}
	return false
}

func isRepeated(s string) bool {
	if s == "" {
		return false
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}
	return true
}

func ValidatePassword(pw string, emailLocalPart string) (bool, string) {
	var reasons []string
	if len(pw) < DefaultStrongPolicy.MinLen {
		reasons = append(reasons, "length ")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, r := range pw {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if DefaultStrongPolicy.RequireUpper && !hasUpper {
		reasons = append(reasons, "upper ")
	}
	if DefaultStrongPolicy.RequireLower && !hasLower {
		reasons = append(reasons, "lower ")
	}
	if DefaultStrongPolicy.RequireNumber && !hasNumber {
		reasons = append(reasons, "number ")
	}
	if DefaultStrongPolicy.RequireSpecial && !hasSpecial {
		reasons = append(reasons, "special ")
	}

	if isRepeated(pw) {
		reasons = append(reasons, "repeated ")
	}
	if isSequential(pw, 4) {
		reasons = append(reasons, "sequential ")
	}

	if DefaultStrongPolicy.DisallowEmailContain && emailLocalPart != "" {
		if strings.Contains(strings.ToLower(pw), strings.ToLower(emailLocalPart)) {
			reasons = append(reasons, "contains_email")
		}
	}

	return len(reasons) == 0, strings.Join(reasons, ",")
}
