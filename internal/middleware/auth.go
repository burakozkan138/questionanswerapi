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
			response := models.NewResponse(false, "Invalid JSON", http.StatusBadRequest, nil, nil)
			response.Write(w)
			return
		}

		userID, err := pkg.ValidateToken(token)
		if err != nil {
			newToken, err := GetAccessFromRefreshToken(r)
			if err != nil {
				response := models.NewResponse(false, "Unauthorized", http.StatusUnauthorized, nil, nil)
				response.Write(w)
				return
			}

			w.Header().Set("Authorization", "Bearer "+newToken)
			userID, err = pkg.ValidateToken(newToken)
			if err != nil {
				response := models.NewResponse(false, "Unauthorized", http.StatusUnauthorized, nil, nil)
				response.Write(w)
				return
			}
		}

		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			response := models.NewResponse(false, "User not found", http.StatusNotFound, nil, nil)
			response.Write(w)
			return
		}

		if user.Blocked {
			response := models.NewResponse(false, "User is blocked", http.StatusForbidden, nil, nil)
			response.Write(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, models.USER_CTX_KEY, userID)
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
