package utils

import "golang.org/x/crypto/bcrypt"

/**
利用hash加密，实现类似laravel的bcrypt()/Hash::make()
 */
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

/**
利用hash验证，实现类似larael Hash::check()
 */
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
