package middleware

import (
	"E-Meeting/configs"
	"E-Meeting/pkg/utils"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return utils.JSONErrorResponse(c, http.StatusUnauthorized, "missing Authorization header")
		}
		log.Printf("Auth Header: %s\n", authHeader)
		// Ensure the token has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.JSONErrorResponse(c, http.StatusUnauthorized, "invalid Authorization header format")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Token String: %s\n", tokenString)
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, utils.JSONErrorResponse(c, http.StatusUnauthorized, "invalid signing method")
			}
			if len(configs.Token.JWTSecret) == 0 {
				return nil, utils.JSONErrorResponse(c, http.StatusInternalServerError, "JWT secret is not configured")
			}
			return []byte(configs.Token.JWTSecret), nil
		})

		if err != nil {
			return utils.JSONErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
		}
		log.Printf("Token: %+v\n", token)
		if !token.Valid {
			return utils.JSONErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			return utils.JSONErrorResponse(c, http.StatusUnauthorized, "unable to parse claims")
		}
		c.Set("user", claims)

		return next(c)
	}
}

func IsAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user")
		log.Printf("User Token: %+v\n", userToken)
		if userToken == nil {
			return utils.JSONErrorResponse(c, http.StatusUnauthorized, "missing or invalid token")
		}

		claims, ok := userToken.(*jwt.MapClaims)
		if !ok {
			return utils.JSONErrorResponse(c, http.StatusUnauthorized, "could not parse token claims")
		}

		isAdmin, ok := (*claims)["is_admin"].(bool)
		if !ok || !isAdmin {
			return utils.JSONErrorResponse(c, http.StatusForbidden, "access denied: admin privileges required")
		}

		return next(c)
	}

}
