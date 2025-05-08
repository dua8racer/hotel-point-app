package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hotel-point-app/internal/models"
	"hotel-point-app/pkg/utils"
)

// AdminOnly adalah middleware untuk memastikan hanya admin yang dapat mengakses
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil user dari context yang sudah diset di middleware auth
		user, exists := c.Get("user")
		if !exists {
			utils.SendUnauthorizedResponse(c)
			c.Abort()
			return
		}

		// Cast ke model User
		userObj, ok := user.(*models.User)
		if !ok {
			utils.SendServerErrorResponse(c, nil)
			c.Abort()
			return
		}

		// Check if user has admin role
		if userObj.Role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, utils.Response{
				Status: "error",
				Error:  "Admin access required",
			})
			c.Abort()
			return
		}

		// User is admin, proceed to next handler
		c.Next()
	}
}
