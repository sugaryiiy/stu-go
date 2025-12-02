package auth

import "github.com/gin-gonic/gin"

// Handler exposes authentication endpoints.
type Handler struct {
	Service Service
}

// RegisterRoutes sets up auth endpoints under the provided route group.
func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Placeholder for future auth routes.
	_ = router
}
