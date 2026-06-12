package ws

import (
	"strings"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/web/result"
)

// AuthResult holds the result of WebSocket authentication.
type AuthResult struct {
	UserID   string
	UserType enums.LoginTypeEnum
	OK       bool
}

// AuthenticateFromToken extracts the user from a query-param token.
// This is the unified auth helper for all WebSocket endpoints.
// Returns an AuthResult; if !OK, an error response has already been written.
func AuthenticateFromToken(c *gin.Context, loginType enums.LoginTypeEnum) AuthResult {
	token := getWebSocketToken(c, loginType)
	if token == "" {
		result.Failure(c, "缺少token", 401)
		c.Abort()
		return AuthResult{OK: false}
	}

	var userID string
	if loginType == enums.LoginTypeConsumer {
		userID = auth.Consumer.GetLoginIDByToken(token)
	} else {
		userID = auth.GetLoginIDByToken(token)
	}

	if userID == "" {
		result.Failure(c, "token无效或已过期", 401)
		c.Abort()
		return AuthResult{OK: false}
	}

	return AuthResult{
		UserID:   userID,
		UserType: loginType,
		OK:       true,
	}
}

func getWebSocketToken(c *gin.Context, loginType enums.LoginTypeEnum) string {
	if c == nil {
		return ""
	}

	if token := strings.TrimSpace(c.Query("token")); token != "" {
		return token
	}

	tokenName := auth.GetTokenName()
	if loginType == enums.LoginTypeConsumer {
		tokenName = auth.Consumer.GetTokenName()
	}
	if token := strings.TrimSpace(c.GetHeader(tokenName)); token != "" {
		return token
	}

	// Some websocket clients can only send auth data via subprotocol/header.
	if protocol := strings.TrimSpace(c.GetHeader("Sec-WebSocket-Protocol")); protocol != "" {
		for _, part := range strings.Split(protocol, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if strings.HasPrefix(strings.ToLower(part), "token.") {
				return strings.TrimSpace(part[len("token."):])
			}
		}
	}

	if authz := strings.TrimSpace(c.GetHeader("Authorization")); authz != "" {
		return authz
	}

	return ""
}
