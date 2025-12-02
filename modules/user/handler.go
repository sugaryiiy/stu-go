package user

import "github.com/gin-gonic/gin"

// Handler wires HTTP routes to the user service.
type Handler struct {
	Service Service
}

// RegisterRoutes sets up user-related routes on the given router group.
func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Placeholder for future user endpoints.
	_ = router
}
