package interfaces

import (
	"fmt"
	"net/http"

	"github.com/KDF5000/nomo/application"
	"github.com/KDF5000/nomo/interfaces/common"
	"github.com/KDF5000/nomo/interfaces/proto"
	"github.com/gin-gonic/gin"
)

// message handler for openapi
type messageHandler struct {
	app application.IMessageHandleApp
}

func NewMessageHandler(app application.IMessageHandleApp) *messageHandler {
	return &messageHandler{app: app}
}

func (h *messageHandler) HandleMessage(c *gin.Context) {
	var request proto.Message
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResonse{
			Code:    common.ErrInvalidParam,
			Message: err.Error(),
		})
		return
	}

	if t, ok := proto.IsValidType(request.Type); !ok {
		c.JSON(http.StatusBadRequest, common.APIResonse{
			Code:    common.ErrInvalidParam,
			Message: fmt.Sprintf("type invalid, must be one of '%s'", t),
		})
		return
	}

	if t, ok := proto.IsValidTheme(request.Theme); !ok {
		c.JSON(http.StatusBadRequest, common.APIResonse{
			Code:    common.ErrInvalidParam,
			Message: fmt.Sprintf("theme invalid, must be one of '%s'", t),
		})
		return
	}

	if (request.Type == "docx" && request.DocxInfo == nil) ||
		(request.Type == "notion" && request.NotionInfo == nil) {
		c.JSON(http.StatusBadRequest, common.APIResonse{
			Code:    common.ErrInvalidParam,
			Message: "type and info mismatched",
		})
		return
	}

	if request.Type == "docx" && !request.DocxInfo.IsValid() {
		c.JSON(http.StatusBadRequest, common.APIResonse{
			Code:    common.ErrInvalidParam,
			Message: fmt.Sprintf("docx info invalid, some field is empty, %+v", request.DocxInfo),
		})
		return
	}

	if request.Type == "notion" && !request.NotionInfo.IsValid() {
		c.JSON(http.StatusBadRequest, common.APIResonse{
			Code:    common.ErrInvalidParam,
			Message: fmt.Sprintf("notion info invalid, some field is empty, %+v", request.NotionInfo),
		})
		return
	}

	if request.Content == "" {
		c.JSON(http.StatusBadRequest, common.APIResonse{
			Code:    common.ErrInvalidParam,
			Message: "content is empty",
		})
		return
	}

	err := h.app.ProcessMessage(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResonse{
			Code:    common.ErrInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResonse{
		Code:    common.ErrSucc,
		Message: "succ",
	})
}
