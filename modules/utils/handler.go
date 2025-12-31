package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type Handler struct {
	service service
}

func NewHandler(db *xorm.Engine) *Handler {
	return &Handler{service: service{repo: repository{db}}}
}
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/getsign", h.getSign)
}

func (h *Handler) getSign(c *gin.Context) {
	util := &SignUtil{}
	err := c.ShouldBindJSON(util)
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = h.service.getSign(util)
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	marshal, err := json.Marshal(util)
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, string(marshal))
}
