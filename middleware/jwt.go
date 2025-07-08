package middleware

import (
	"case_study_api/config"
	"case_study_api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildErrorResponse("missing or malformed JWT"))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.App.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildErrorResponse("invalid or expired JWT"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildErrorResponse("invalid JWT claims"))
			return
		}

		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("user_role", claims["role"].(string))
		c.Next()
	}
}

func RoleAuth(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.BuildErrorResponse("role not found in token"))
			return
		}

		userRole := roleVal.(string)
		for _, role := range requiredRoles {
			if role == userRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, utils.BuildErrorResponse("access forbidden for role: "+userRole))
	}
}
