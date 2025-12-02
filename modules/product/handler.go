package product

import "github.com/gin-gonic/gin"

// Handler exposes product endpoints.
type Handler struct {
	Service Service
}

// RegisterRoutes sets up product endpoints under the provided route group.
func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Placeholder for future product routes.
	_ = router
}
