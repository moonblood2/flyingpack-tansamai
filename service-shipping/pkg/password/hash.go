package password

import "golang.org/x/crypto/bcrypt"

// BcryptHash return hash of password that use Bcrypt Algorithm.
func BcryptHash(password string) (string, error) {
	hashAndSaltByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashAndSaltByte), nil
}
