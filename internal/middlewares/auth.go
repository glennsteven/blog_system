package middlewares

import (
	"blog-system/internal/config"
	"blog-system/internal/entities"
	"blog-system/internal/helper"
	"blog-system/internal/resources"
	"blog-system/pkg/viper"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response resources.Response
		v, err := viper.GlobalConfig()
		if err != nil {
			return
		}
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					return []byte(v.GetString("JWT_KEY")), nil
				})

				if err != nil {
					response.Code = http.StatusUnauthorized
					response.Message = "unauthorized"
					helper.ResponseJSON(w, http.StatusUnauthorized, response)
					return
				}

				if !token.Valid {
					response.Code = http.StatusUnauthorized
					response.Message = "unauthorized"
					helper.ResponseJSON(w, http.StatusUnauthorized, response)
					return
				}
				next.ServeHTTP(w, r)
			}
		} else {
			response.Code = http.StatusUnauthorized
			response.Message = "authorization header is required"
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}
	})
}

func ParseJWT(cfg config.Jwt, tokenString string) (jwt.MapClaims, error) {
	bearerToken := strings.Split(tokenString, " ")
	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.Key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func DecodeJWT(cfg config.Jwt, tokenString string) (*entities.UserClaim, error) {
	claims, err := ParseJWT(cfg, tokenString)
	if err != nil {
		return nil, errors.New("unable to parse JWT token")
	}

	fullName, ok := claims["FullName"].(string)
	if !ok {
		return nil, errors.New("FullName field missing in claims")
	}

	userIdFloat64, ok := claims["Id"].(float64)
	if !ok {
		return nil, errors.New("id field missing in claims")
	}

	roleIdFloat64, ok := claims["RoleId"].(float64)
	if !ok {
		return nil, errors.New("id field missing in claims")
	}
	userId := int(userIdFloat64)
	roleId := int64(roleIdFloat64)

	user := &entities.UserClaim{
		Id:       userId,
		FullName: fullName,
		RoleId:   roleId,
	}

	return user, nil
}
