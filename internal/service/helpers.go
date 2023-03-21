package service

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func GenSalt(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ``
	}
	return hex.EncodeToString(b)
}

func GenHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func CheckHash(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func LuhnCheck(order string) bool {
	return true
}
