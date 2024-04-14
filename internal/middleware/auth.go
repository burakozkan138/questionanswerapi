package middleware

import (
	"burakozkan138/questionanswerapi/internal/database"
	"burakozkan138/questionanswerapi/internal/models"
	"burakozkan138/questionanswerapi/pkg"
	"context"
	"errors"
	"net/http"
	"strings"
)

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			db   = database.GetDB()
			user models.User
		)

		token, err := GetTokenFromHeader(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := pkg.ValidateToken(token)
		if err != nil {
			newToken, err := GetAccessFromRefreshToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			w.Header().Set("Authorization", "Bearer "+newToken)
			userID, err = pkg.ValidateToken(newToken)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		}

		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Blocked {
			http.Error(w, "User is blocked", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetTokenFromHeader(r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")

	if authorization == "" {
		return "", errors.New("authorization header missing")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid authorization token")
	}

	return strings.TrimPrefix(authorization, "Bearer "), nil
}

func GetAccessFromRefreshToken(r *http.Request) (string, error) {
	refreshtoken := r.Header.Get("refresh_token")
	if refreshtoken == "" {
		return "", errors.New("refresh token missing")
	}

	newToken, err := pkg.RefreshToken(refreshtoken)
	if err != nil {
		return "", errors.New("failed to refresh token")
	}

	return newToken, nil
}
