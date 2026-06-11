package middleware

import (
	"net/http"
	"strings"
	"time"

	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gin-gonic/gin"
)

const ClerkClaimsKey = "clerk_claims"

func ClerkAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" || token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"data": nil,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "missing authorization token",
				},
				"meta": gin.H{
					"request_id": c.GetString("request_id"),
					"timestamp":  time.Now().UTC().Format(time.RFC3339),
				},
			})
			return
		}
		claims, err := clerkjwt.Verify(c.Request.Context(), &clerkjwt.VerifyParams{
			Token: token,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"data": nil,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "invalid or expired token",
				},
				"meta": gin.H{
					"request_id": c.GetString("request_id"),
					"timestamp":  time.Now().UTC().Format(time.RFC3339),
				},
			})
			return
		}
		c.Set(ClerkClaimsKey, claims)
		c.Next()
	}
}
