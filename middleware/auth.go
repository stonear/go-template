package middleware

import (
	"errors"
	"net/http"
	"strings"

	"aidanwoods.dev/go-paseto"
	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/response"
	"github.com/stonear/go-template/validator"
	"go.uber.org/zap"
)

func (m *middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			m.log.Error("failed to extract bearer token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusForbidden, response.New(
				response.CodeNotPermitted,
				validator.Message(err),
			))
			return
		}

		parser := paseto.NewParser()
		parsedToken, err := parser.ParseV4Local(m.key, token, nil)
		if err != nil {
			m.log.Error("failed to parse paseto token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusForbidden, response.New(
				response.CodeNotPermitted,
				validator.Message(err),
			))
			return
		}

		id, _ := parsedToken.GetString("id")
		username, _ := parsedToken.GetString("username")
		c.Set("x-id", id)
		c.Set("x-username", username)

		c.Next()
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	return strings.TrimPrefix(header, "Bearer "), nil
}
