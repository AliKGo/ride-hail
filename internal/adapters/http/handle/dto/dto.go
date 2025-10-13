package dto

import (
	"ride-hail/internal/adapters/http/handle/dto/validate"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"strings"
)

func ValidateLogin(email string, password string) (bool, string) {
	var res []string

	if !validate.ValidateEmail(email, true) {
		res = append(res, "email is invalid")
	}

	localPart := ""
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		localPart = parts[0]
	}

	if ok, msg := validate.ValidatePassword(password, localPart); !ok {
		res = append(res, msg)
	}

	if len(res) == 0 {
		return true, ""
	}

	return false, strings.Join(res, ", ")
}

func GetRole(mode string, user *models.User) {
	switch mode {
	case types.ModeAdmin:
		user.Role = types.RoleAdmin
	case types.ModeDAL:
		user.Role = types.RoleDriver
	case types.ModeRide:
		user.Role = types.RoleCustomer
	}
}
