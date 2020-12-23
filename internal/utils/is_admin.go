package utils

import "github.com/xuxusheng/time-frequency-be/internal/model"

func IsAdmin(roles []model.Role) bool {
	for _, role := range roles {
		if role == model.Admin {
			return true
		}
	}
	return false
}
