package interfaces

import (
	"github.com/gin-gonic/gin"

	"github.com/KDF5000/nomo/application"
)

type BindInfoHandler struct {
	bindApp application.IBindInfoApp
}

func NewBindHandler(app application.IBindInfoApp) *BindInfoHandler {
	return &BindInfoHandler{bindApp: app}
}

func (h *BindInfoHandler) BindLark(c *gin.Context) {
}

func (h *BindInfoHandler) BindWX(c *gin.Context) {
}
