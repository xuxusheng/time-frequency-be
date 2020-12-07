package pwd

import "golang.org/x/crypto/bcrypt"

// 密码 + 盐后计算 hash
func EncodePWD(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 对比密码和数据库中的 hash 值
func ComparePWD(hash, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}
