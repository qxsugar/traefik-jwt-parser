package traefik_jwt_parser

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

type Config struct {
	TokenKey  string
	SecretKey string
	TrustKeys []string
}

func CreateConfig() *Config {
	standardTrustKeys := []string{"aud", "exp", "jti", "iat", "sub", "nbf", "sub"}
	config := &Config{
		TokenKey:  "Authorization",
		SecretKey: "traefik-jwt",
		TrustKeys: []string{},
	}
	config.TrustKeys = append(config.TrustKeys, standardTrustKeys...)
	return config
}

type JWTParser struct {
	context context.Context
	next    http.Handler
	config  *Config
}

func (p JWTParser) findToken(request *http.Request) string {
	token := request.Header.Get(p.config.TokenKey)
	if token == "" {
		return request.URL.Query().Get(p.config.TokenKey)
	}
	return token
}

func (p JWTParser) formatTrustKey(key string) string {
	return strings.ToUpper(fmt.Sprintf("X-%s", key))
}

func (p JWTParser) resetTrustHeaders(request *http.Request) {
	for _, key := range p.config.TrustKeys {
		request.Header.Del(p.formatTrustKey(key))
	}
}

func (p JWTParser) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	p.resetTrustHeaders(request)

	jwtToken := p.findToken(request)
	if jwtToken == "" {
		p.next.ServeHTTP(writer, request)
		return
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) { return []byte(p.config.SecretKey), nil })
	if err != nil {
		p.next.ServeHTTP(writer, request)
		return
	}

	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		for _, trustKey := range p.config.TrustKeys {
			if value, ok := mapClaims[trustKey]; ok {
				v := fmt.Sprintf("%v", value)
				request.Header.Set(p.formatTrustKey(trustKey), v)
			}
		}
	}

	p.next.ServeHTTP(writer, request)
	return
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &JWTParser{
		context: ctx,
		next:    next,
		config:  config,
	}, nil
}

var _ http.Handler = &JWTParser{}
