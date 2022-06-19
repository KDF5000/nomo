package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/KDF5000/nomo/application"
	"github.com/KDF5000/nomo/infrastructure/utils"
	"github.com/KDF5000/nomo/interfaces/proto"
	itemplate "github.com/KDF5000/nomo/interfaces/template"
	"github.com/KDF5000/pkg/log"
)

type posterHandler struct {
	posterApp application.IPosterApp
}

func NewPosterHandler(app application.IPosterApp) *posterHandler {
	return &posterHandler{posterApp: app}
}

func (h *posterHandler) GenPoster(c *gin.Context) {
	var request proto.PosterRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if !request.IsValid() {
		log.Errorf("invalid request %+v", request)
		c.JSON(http.StatusBadRequest, "invalid request parameter")
		return
	}

	viewData := itemplate.TPLMemoViewData{
		UserName:  request.UserName,
		CreatedAt: request.CreatedAt,
	}

	elements := utils.ScanContent(request.Content)
	for _, ele := range elements {
		viewData.ContentElements = append(viewData.ContentElements, itemplate.ContntElement{
			IsTag:   ele.IsTag,
			Content: ele.Text,
		})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	log.Infof("%v", viewData)
	data, err := h.posterApp.GenPoster(c.Request.Context(), uint(id), viewData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.WriteString(string(data))
}

func (h *posterHandler) Screenshot(c *gin.Context) {
	var request proto.ScreenshotRequst
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if !request.IsValidUrl() {
		c.JSON(http.StatusBadRequest, "invalid url, must be like http://xxx.com or https://xxx.com")
		return
	}
	data, err := h.posterApp.Screnshot(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.WriteString(string(data))
}
