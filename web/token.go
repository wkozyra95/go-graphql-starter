package web

import (
	"context"
	"crypto/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	conf "github.com/wkozyra95/go-graphql-starter/config"
	"github.com/wkozyra95/go-graphql-starter/model/mongo"
	"github.com/wkozyra95/go-graphql-starter/web/schema"
	"gopkg.in/mgo.v2/bson"
)

type jwtProvider struct {
	jwtKey []byte
	header string
}

func newJwtProvider(config conf.Config) (jwtProvider, error) {
	log.Info("Create jwtProvider")
	const keySize = 64
	jwtKey := make([]byte, keySize)
	_, err := rand.Read(jwtKey)
	if err != nil {
		return jwtProvider{}, err
	}
	return jwtProvider{
		jwtKey: []byte("ewjfhbweruhf"),
		header: "X-Auth-Token",
	}, nil
}

func (jp jwtProvider) GenerateToken(id bson.ObjectId) string {
	log.Debugf("GenerateToken [%s]", id.Hex())
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 2400).Unix()

	signedToken, signErr := token.SignedString(jp.jwtKey)
	if signErr != nil {
		log.Error("[ASSERT] Unable to sign token")
		return ""
	}
	return signedToken
}

func (jp jwtProvider) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get(jp.header)
		if token == "" {
			log.Info("Missing auth token")
			next.ServeHTTP(w, r)
			return
		}

		parsed, parseErr := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
			return jp.jwtKey, nil
		})
		if validationErr, ok := parseErr.(*jwt.ValidationError); ok &&
			validationErr.Errors&jwt.ValidationErrorExpired != 0 {
			log.Warnf("token expired [%s]", parseErr.Error())
			next.ServeHTTP(w, r)
			return
		}
		if parseErr != nil {
			log.Warnf("Unable to parse token [%s]", parseErr.Error())
			next.ServeHTTP(w, r)
			return
		}

		claims, assertTypeOk := parsed.Claims.(jwt.MapClaims)
		if !parsed.Valid || !assertTypeOk {
			log.Warn("Token is not valid")
			next.ServeHTTP(w, r)
			return
		}

		converted, convertErr := mongo.ConvertToObjectID(claims["id"].(string))
		if convertErr != nil {
			next.ServeHTTP(w, r)
			return
		}

		idKey := schema.CurrentUserKey
		idVal := converted

		ctx := context.WithValue(r.Context(), idKey, idVal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
