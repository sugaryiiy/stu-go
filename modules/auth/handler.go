package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler exposes authentication endpoints.
type Handler struct {
	Service Service
}

// RegisterRoutes sets up auth endpoints under the provided route group.
func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", h.Login)
	router.POST("/refresh", h.Refresh)
	router.GET("/me", h.requireAuth(), h.Me)
}

// Login authenticates a user and returns a new token pair.
func (h Handler) Login(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.Service.Login(creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Refresh exchanges a refresh token for a new pair of tokens.
func (h Handler) Refresh(c *gin.Context) {
	var payload struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.Service.Refresh(payload.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Me returns the authenticated user's claims.
func (h Handler) Me(c *gin.Context) {
	claims, _ := c.Get("claims")
	c.JSON(http.StatusOK, claims)
}

func (h Handler) requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		claims, err := h.Service.ValidateAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("claims", claims)
	}
}
