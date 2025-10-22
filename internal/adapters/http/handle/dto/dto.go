package dto

import (
	"ride-hail/internal/adapters/http/handle/dto/validate"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"strings"
)

var mode string

func ValidateLogin(u *models.User) (bool, string) {
	var res []string

	if !validate.ValidateEmail(u.Email, true) {
		res = append(res, "email is invalid")
	}

	localPart := ""
	parts := strings.Split(u.Email, "@")
	if len(parts) == 2 {
		localPart = parts[0]
	}

	if ok, msg := validate.ValidatePassword(u.Password, localPart); !ok {
		res = append(res, msg)
	}

	if len(res) == 0 {
		getRole(u)
		return true, ""
	}

	return false, strings.Join(res, ", ")
}

func InitMode(m string) {
	mode = m
}

func getRole(user *models.User) {
	switch mode {
	case types.ModeAdmin:
		user.Role = types.RoleAdmin
	case types.ModeDAL:
		user.Role = types.RoleDriver
	case types.ModeRide:
		user.Role = types.RoleCustomer
	}
}
