package util

import (
	"encoding/hex"

	"github.com/minio/highwayhash"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPwd(pwd *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*pwd = string(hash)
	return nil
}

func ValidatePwd(hash, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)); err!= nil {
		return false
	}
	return true
}

func Hash(s string) string {
	var hash = highwayhash.Sum128([]byte(s), make([]byte, 32))
	return hex.EncodeToString(hash[:])
}

// func SendCode(email, code string) {
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", "luckycatxk@gmail.com")
// 	m.SetHeader("To", email)
// 	m.SetAddressHeader("Cc", "John@example.com", "John")
// 	m.SetHeader("Subject", "Code")
// 	m.SetBody("text/plain", "Your code is: " + code)

// 	d := gomail.NewDialer("smtp.google.com", 587, "luckycatxk@gmail.com", mailpwd)

// 	if err := d.DialAndSend(m); err != nil {
// 		panic(err)
// 	}
// }
