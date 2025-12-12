package user

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type Handler struct {
	service service
}

func NewHandler(db *xorm.Engine) *Handler {
	service1 := service{repo: repository{db}}
	return &Handler{service: service1}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/addUser", h.Create)
	router.GET("getUser/:id", h.GetByID)
	router.GET("/getUserList", h.List)
	router.POST("/login")
}

func (h *Handler) Create(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := h.service.Create(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) List(c *gin.Context) {
	users, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
