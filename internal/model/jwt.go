package model

type JWTClaims struct {
	UID   uint   `json:"uid"` // 用户 ID
	Roles []Role `json:"roles"`
}
