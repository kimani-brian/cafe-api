package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole takes a list of allowed roles. If the user isn't on the list, they get blocked.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Grab the role that RequireAuth() saved in the context
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			return
		}

		// Check if their role is in the allowed list
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: insufficient permissions"})
			return
		}

		// If they have the right role, let them through
		c.Next()
	}
}
