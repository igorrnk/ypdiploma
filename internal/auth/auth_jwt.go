package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"net/http"
	"time"
)

var hmacSampleSecret = []byte("SecretKey")

type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthenticatorJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authString := r.Header.Get("Authorization")
		if len(authString) < 7 || authString[0:6] != "Bearer" {
			http.Error(w, "Wrong token type", http.StatusUnauthorized)
			return
		}
		tokenString := authString[7:]

		user, err := UserByToken(tokenString)
		if errors.Is(err, ErrWrongToken) {
			http.Error(w, "Wrong token", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, "Wrong token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "User", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TokenByUser(user *model.User) string {
	claims := &UserClaims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(hmacSampleSecret)
	return tokenString
}

func UserByToken(tokenString string) (*model.User, error) {
	user := model.User{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrWrongToken
			}
			return hmacSampleSecret, nil
		})
	if err != nil {
		return nil, ErrWrongToken
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		user.ID = claims.UserID
	}
	return &user, nil
}
