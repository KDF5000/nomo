package interfaces

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KDF5000/nomo/application"
	"github.com/KDF5000/nomo/infrastructure/message/lark_message"
	"github.com/KDF5000/nomo/interfaces/common"
	"github.com/KDF5000/pkg/log"
)

type larkMessageHandler struct {
	messageHandleApp application.ILarkMessageHandleApp
}

func NewLarkMessageHandler(app application.ILarkMessageHandleApp) *larkMessageHandler {
	return &larkMessageHandler{messageHandleApp: app}
}

func (h *larkMessageHandler) UrlVerification(c *gin.Context) {
	var event lark_message.UrlVerificationEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.messageHandleApp.VerifyURL(context.TODO(), &event)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *larkMessageHandler) HandleMessage(c *gin.Context) {
	if c.Request == nil || c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// log.Infof("data ==> %+v", string(data))
	var event lark_message.LarkMessageEvent
	if err := json.Unmarshal(data, &event); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if event.Schema == "" {
		var verifyEvent lark_message.UrlVerificationEvent
		if err := json.Unmarshal(data, &verifyEvent); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		res, err := h.messageHandleApp.VerifyURL(context.TODO(), &verifyEvent)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, res)
		return
	}

	// log.Infof("%+v", event)
	go func() {
		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		defer cancel()
		if err := h.messageHandleApp.ProcessMessage(ctx, &event); err != nil {
			log.Error("failed to process lark message", log.String("err", err.Error()))
		}
	}()

	c.JSON(http.StatusOK, common.APIResonse{
		Code:    0,
		Message: "succ",
	})
}
