package pkg

import (
	"burakozkan138/questionanswerapi/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateToken(userID uuid.UUID, exp string, otherClaims map[string]interface{}) (string, error) {
	config := config.JwtConfig

	expDuration, err := time.ParseDuration(exp)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(expDuration).Unix(),
		"iat":    time.Now().Unix(),
		"iss":    config.Issuer,
		"aud":    config.Audience,
	}

	for key, value := range otherClaims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (uuid.UUID, error) {
	config := config.JwtConfig

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return uuid.Max, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Max, err
	}

	userID, err := uuid.Parse(claims["userId"].(string))
	if err != nil {
		return uuid.Max, err
	}

	return userID, nil
}

func CreateAccessToken(userID uuid.UUID) (string, error) {
	return CreateToken(userID, config.JwtConfig.AccessExpiresIn, nil)
}

func CreateRefreshToken(userID uuid.UUID) (string, error) {
	return CreateToken(userID, config.JwtConfig.RefreshExpiresIn, nil)
}

func RefreshToken(refreshToken string) (string, error) {
	userID, err := ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err := CreateToken(userID, config.JwtConfig.AccessExpiresIn, nil)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
