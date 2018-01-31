package web

import (
	"context"
	"crypto/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	conf "github.com/wkozyra95/go-graphql-starter/config"
	"github.com/wkozyra95/go-graphql-starter/errors"
	"gopkg.in/mgo.v2/bson"
)

type contextKeyType string

const contextUserID contextKeyType = "userId"

type jwtProvider struct {
	jwtKey []byte
	header string
}

func newJwtProvider(config conf.Config) (jwtProvider, error) {
	const keySize = 64
	jwtKey := make([]byte, keySize)
	_, err := rand.Read(jwtKey)
	if err != nil {
		return jwtProvider{}, err
	}
	return jwtProvider{
		jwtKey: jwtKey,
		header: "X-Auth-Token",
	}, nil
}

func (jp jwtProvider) generate(id bson.ObjectId) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return token.SignedString(jp.jwtKey)
}

func (jp jwtProvider) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get(jp.header)
		if token == "" {
			log.Info("Missing auth token")
			_ = writeJSONResponse(w, http.StatusUnauthorized, map[string]string{
				"reason": errors.ErrNotLoggedIn,
			})
			return
		}

		parsed, parseErr := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
			return jp.jwtKey, nil
		})
		if validationErr, ok := parseErr.(*jwt.ValidationError); ok &&
			validationErr.Errors&jwt.ValidationErrorExpired != 0 {
			log.Warn("token expired")
			_ = writeJSONResponse(w, http.StatusUnauthorized, map[string]string{
				"reason": errors.ErrExpired,
			})
			return
		}
		if parseErr != nil {
			log.Warn("Unable to parse token")
			_ = writeJSONResponse(w, http.StatusUnauthorized, map[string]string{
				"reason": errors.ErrMalformed,
			})
			return
		}

		claims, assertTypeOk := parsed.Claims.(jwt.MapClaims)
		if !parsed.Valid || !assertTypeOk {
			log.Warn("Token is not valid")
			_ = writeJSONResponse(w, http.StatusUnauthorized, map[string]string{
				"reason": errors.ErrMalformed,
			})
			return
		}

		idKey := contextUserID
		idVal := claims["id"]

		ctx := context.WithValue(r.Context(), idKey, idVal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
